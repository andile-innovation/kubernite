package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	kubernetesRestClient "k8s.io/client-go/rest"
	"kubernite/pkg/git"
	kubernetesManifest "kubernite/pkg/kubernetes/manifest"
	"os"
)

func main() {

	gitRepo, err := git.NewRepositoryFromFilePath(os.Getenv("PLUGIN_DEPLOYMENT_TAG_REPOSITORY_PATH"))
	if err != nil {
		log.Fatal(err)
	}

	latestTag, err := gitRepo.GetLatestTagName()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("latest tag! ", latestTag)
	updateDeploymentFile(latestTag)

	var config = new(kubernetesRestClient.Config)
	config.Host = os.Getenv("PLUGIN_KUBERNETES_SERVER")
	config.TLSClientConfig.CAData = []byte(os.Getenv("PLUGIN_KUBERNETES_CERT_DATA"))
	config.TLSClientConfig.CertData = []byte(os.Getenv("PLUGIN_KUBERNETES_CLIENT_CERT_DATA"))
	config.TLSClientConfig.KeyData = []byte(os.Getenv("PLUGIN_KUBERNETES_CLIENT_KEY_DATA"))

	// create the client set
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	for {
		pods, err := clientset.CoreV1().Pods("dev").List(v1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster!\n", len(pods.Items))
		break
	}
}

func updateDeploymentFile(tag string) {
	deploymentFile, err := kubernetesManifest.NewDeploymentFromFile(os.Getenv("PLUGIN_KUBERNETES_DEPLOYMENT_FILE_PATH"))
	if err != nil {
		log.Fatal(err)
	}

	if err := deploymentFile.UpdatePodTemplateAnnotations(
		"kubernetes.io/change-cause",
		"Update image to new version",
	); err != nil {
		log.Fatal(err)
	}

	if err := deploymentFile.WriteAtPath("output.yaml"); err != nil {
		log.Fatal(err)
	}
}

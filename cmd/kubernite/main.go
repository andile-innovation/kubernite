package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	kubernetesRestClient "k8s.io/client-go/rest"
	kuberniteConfig "kubernite/configs/kubernite"
	"kubernite/pkg/git"
	kubernetesManifest "kubernite/pkg/kubernetes/manifest"
	"os"
	"strings"
)

func main() {

	conf, err := kuberniteConfig.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(conf.DeploymentRepositoryPath)

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if len(pair) < 2 {
			continue
		}
		if strings.HasPrefix(pair[0], "DRONE") || strings.HasPrefix(pair[0], "CI") {
			fmt.Printf("%s = %s\n", pair[0], pair[1])
		}
	}

	// open git repository
	gitRepo, err := git.NewRepositoryFromFilePath(conf.DeploymentRepositoryPath)
	if err != nil {
		log.Fatal(err)
	}

	// get the latest git tag on the repository
	latestTag, err := gitRepo.GetLatestTagName()
	if err != nil {
		log.Fatal(err)
	}

	// open deployment file
	deploymentFile, err := kubernetesManifest.NewDeploymentFromFile(conf.KubernetesDeploymentFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// update deployment file annotations with tag
	if err := deploymentFile.UpdatePodTemplateAnnotations(
		"kubernetes.io/change-cause",
		fmt.Sprintf("update image to %s", latestTag),
	); err != nil {
		log.Fatal(err)
	}

	var config = new(kubernetesRestClient.Config)
	config.Host = conf.KubernetesServer
	config.TLSClientConfig.CAData = []byte(conf.KubernetesCertData)
	config.TLSClientConfig.CertData = []byte(conf.KubernetesClientCertData)
	config.TLSClientConfig.KeyData = []byte(conf.KubernetesClientKeyData)

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

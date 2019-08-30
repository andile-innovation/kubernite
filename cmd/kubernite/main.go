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
)

func main() {
	// parse configuration
	kuberniteConf, err := kuberniteConfig.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	// handle build event
	switch kuberniteConf.BuildEvent {
	case git.TagEvent:
		if err := handleTagEvent(kuberniteConf); err != nil {
			log.Fatal(err)
		}
	default:
		if err := handleOtherEvent(kuberniteConf); err != nil {
			log.Fatal(err)
		}
	}

	var config = new(kubernetesRestClient.Config)
	config.Host = kuberniteConf.KubernetesServer
	config.TLSClientConfig.CAData = []byte(kuberniteConf.KubernetesCertData)
	config.TLSClientConfig.CertData = []byte(kuberniteConf.KubernetesClientCertData)
	config.TLSClientConfig.KeyData = []byte(kuberniteConf.KubernetesClientKeyData)

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

func handleTagEvent(kuberniteConf *kuberniteConfig.Config) error {
	// open git repository
	gitRepo, err := git.NewRepositoryFromFilePath(kuberniteConf.DeploymentRepositoryPath)
	if err != nil {
		log.Fatal(err)
	}

	// get the latest git tag on the repository
	latestTag, err := gitRepo.GetLatestTagName()
	if err != nil {
		log.Fatal(err)
	}

	// open deployment file
	deploymentFile, err := kubernetesManifest.NewDeploymentFromFile(kuberniteConf.KubernetesDeploymentFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// update deployment file annotations with tag and event information
	if err := deploymentFile.UpdateAnnotations(
		"kubernetes.io/change-cause",
		fmt.Sprintf("kubernite handled tag event - image updated to %s", latestTag),
	); err != nil {
		log.Fatal(err)
	}
	if err := deploymentFile.UpdatePodTemplateAnnotations(
		"kubernetes.io/change-cause",
		fmt.Sprintf("kubernite handled tag event - image updated to %s", latestTag),
	); err != nil {
		log.Fatal(err)
	}

	if kuberniteConf.DryRun {
		log.Info("____tag event dry run____")
		log.Info(fmt.Sprintf("kubectl apply -f %s", kuberniteConf.KubernetesDeploymentFilePath))
		return nil
	}

	return nil
}

func handleOtherEvent(kuberniteConf *kuberniteConfig.Config) error {
	// open git repository
	gitRepo, err := git.NewRepositoryFromFilePath(kuberniteConf.DeploymentRepositoryPath)
	if err != nil {
		log.Fatal(err)
	}

	// get the latest commit hash in the repository
	latestCommitHash, err := gitRepo.GetLatestCommitHash()
	if err != nil {
		log.Fatal(err)
	}

	// open deployment file
	deploymentFile, err := kubernetesManifest.NewDeploymentFromFile(kuberniteConf.KubernetesDeploymentFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// update deployment file annotations with tag and event information
	if err := deploymentFile.UpdateAnnotations(
		"kubernetes.io/change-cause",
		fmt.Sprintf("kubernite handled %s event - commit hash %s", kuberniteConf.BuildEvent, latestCommitHash),
	); err != nil {
		log.Fatal(err)
	}
	if err := deploymentFile.UpdatePodTemplateAnnotations(
		"kubernetes.io/change-cause",
		fmt.Sprintf("kubernite handled %s event - commit hash %s", kuberniteConf.BuildEvent, latestCommitHash),
	); err != nil {
		log.Fatal(err)
	}

	if kuberniteConf.DryRun {
		log.Info(fmt.Sprintf("____%s event dry run____", kuberniteConf.BuildEvent))
		log.Infof(fmt.Sprintf("kubectl apply -f %s", kuberniteConf.KubernetesDeploymentFilePath))
		return nil
	}

	return nil
}

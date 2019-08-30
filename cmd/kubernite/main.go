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
	"time"
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
		fmt.Sprintf(
			"kubernite handled tag event @ %s - image updated to %s",
			time.Now().Format("Jan-02-2006 15:04:05"),
			latestTag,
		),
	); err != nil {
		log.Fatal(err)
	}
	if err := deploymentFile.UpdatePodTemplateAnnotations(
		"kubernetes.io/change-cause",
		fmt.Sprintf(
			"kubernite handled tag event @ %s - image updated to %s",
			time.Now().Format("Jan-02-2006 15:04:05"),
			latestTag,
		),
	); err != nil {
		log.Fatal(err)
	}

	if kuberniteConf.DryRun {
		log.Info("____tag event dry run____")
		log.Info(fmt.Sprintf("kubectl apply -f %s", kuberniteConf.KubernetesDeploymentFilePath))
		deploymentFileContents, err := deploymentFile.GetDeploymentFileContents()
		if err != nil {
			log.Fatal(err)
		}
		log.Info(fmt.Sprintf("\n%s", deploymentFileContents))
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
		fmt.Sprintf(
			"kubernite handled %s event @ %s - commit hash %s",
			time.Now().Format("Jan-02-2006 15:04:05"),
			kuberniteConf.BuildEvent,
			latestCommitHash,
		),
	); err != nil {
		log.Fatal(err)
	}
	if err := deploymentFile.UpdatePodTemplateAnnotations(
		"kubernetes.io/change-cause",
		fmt.Sprintf(
			"kubernite handled %s event @ %s - commit hash %s",
			kuberniteConf.BuildEvent,
			time.Now().Format("Jan-02-2006 15:04:05"),
			latestCommitHash,
		),
	); err != nil {
		log.Fatal(err)
	}

	// write the deployment file

	if kuberniteConf.DryRun {
		log.Info(fmt.Sprintf("____%s event dry run____", kuberniteConf.BuildEvent))
		log.Info(fmt.Sprintf("kubectl apply -f %s", kuberniteConf.KubernetesDeploymentFilePath))
		deploymentFileContents, err := deploymentFile.GetDeploymentFileContents()
		if err != nil {
			log.Fatal(err)
		}
		log.Info(fmt.Sprintf("\n%s", deploymentFileContents))
		return nil
	}

	return nil
}

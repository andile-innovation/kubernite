package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	kuberniteConfig "kubernite/configs/kubernite"
	kubernetesClient "kubernite/internal/pkg/kubernetes/client"
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
	var deploymentFile *kubernetesManifest.Deployment
	switch kuberniteConf.BuildEvent {
	case git.TagEvent:
		if deploymentFile, err = handleTagEvent(kuberniteConf); err != nil {
			log.Fatal(err)
		}
	default:
		if deploymentFile, err = handleOtherEvent(kuberniteConf); err != nil {
			log.Fatal(err)
		}
	}

	if kuberniteConf.DryRun {
		log.Info(fmt.Sprintf("____%s event dry run____", kuberniteConf.BuildEvent))
		log.Info(fmt.Sprintf("kubectl apply -f %s", kuberniteConf.KubernetesDeploymentFilePath))
		deploymentFileContents, err := deploymentFile.GetDeploymentFileContents()
		if err != nil {
			log.Fatal(err)
		}
		log.Info(fmt.Sprintf("\n%s", deploymentFileContents))
		return
	}

	// write file
	if err := deploymentFile.WriteAtPath("output.yaml"); err != nil {
		log.Fatal(err)
	}

	// create a kubernetes client
	kubeClient, err := kubernetesClient.NewClientFromKuberniteConfig(kuberniteConf)
	if err != nil {
		log.Fatal(err)
	}

	// apply updated deployment file
	deploymentFileJSON, err := deploymentFile.ToJSON()
	if err != nil {
		log.Fatal(err.Error())
	}

	response := kubeClient.RESTClient().Verb("apply").Body(deploymentFileJSON).Do()
	if response.Error() != nil {
		log.Fatal(response.Error())
	}
	responseString, err := response.Raw()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(string(responseString))
}

func handleTagEvent(kuberniteConf *kuberniteConfig.Config) (*kubernetesManifest.Deployment, error) {
	// open git repository
	gitRepo, err := git.NewRepositoryFromFilePath(kuberniteConf.DeploymentRepositoryPath)
	if err != nil {
		return nil, err
	}

	// get the latest git tag on the repository
	latestTag, err := gitRepo.GetLatestTagName()
	if err != nil {
		return nil, err
	}

	// open deployment file
	deploymentFile, err := kubernetesManifest.NewDeploymentFromFile(kuberniteConf.KubernetesDeploymentFilePath)
	if err != nil {
		return nil, err
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

	return deploymentFile, nil
}

func handleOtherEvent(kuberniteConf *kuberniteConfig.Config) (*kubernetesManifest.Deployment, error) {
	// open git repository
	gitRepo, err := git.NewRepositoryFromFilePath(kuberniteConf.DeploymentRepositoryPath)
	if err != nil {
		log.Fatal(err)
	}

	// get the latest commit hash in the repository
	latestCommitHash, err := gitRepo.GetLatestCommitHash()
	if err != nil {
		return nil, err
	}

	// open deployment file
	deploymentFile, err := kubernetesManifest.NewDeploymentFromFile(kuberniteConf.KubernetesDeploymentFilePath)
	if err != nil {
		return nil, err
	}

	// update deployment file annotations with tag and event information
	if err := deploymentFile.UpdateAnnotations(
		"kubernetes.io/change-cause",
		fmt.Sprintf(
			"kubernite handled %s event @ %s - commit hash %s",
			kuberniteConf.BuildEvent,
			time.Now().Format("Jan-02-2006 15:04:05"),
			latestCommitHash,
		),
	); err != nil {
		return nil, err
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
		return nil, err
	}

	return deploymentFile, nil
}

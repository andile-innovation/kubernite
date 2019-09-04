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
	deploymentFile, err := handleDeployment(kuberniteConf)
	if err != nil {
		log.Fatal(err)
	}

	// if this is a dry run, print out deployment file to be updated
	if kuberniteConf.DryRun {
		log.Info(fmt.Sprintf("____%s event dry run____", kuberniteConf.BuildEvent))
		log.Info(fmt.Sprintf("kubectl apply -f %s", kuberniteConf.DeploymentFilePath))
		log.Info(fmt.Sprintf("\n%s", deploymentFile.String()))
		return
	}

	// create a kubernetes client
	kubeClient, err := kubernetesClient.NewClientFromKuberniteConfig(kuberniteConf)
	if err != nil {
		log.Fatal(err)
	}

	// apply the deployment
	deploymentClient := kubeClient.Clientset.AppsV1().Deployments(deploymentFile.Namespace)
	if _, err := deploymentClient.Update(deploymentFile.Deployment); err != nil {
		log.Fatal(err)
	}

	// write file
	if err := deploymentFile.WriteToYAMLAtPath(kuberniteConf.DeploymentFilePath); err != nil {
		log.Fatal(err)
	}

	//commit deployment if set
	if kuberniteConf.CommitDeployment {
		if err := commitDeployment(kuberniteConf); err != nil {
			log.Fatal(err)
		}
	}
}

func commitDeployment(kuberniteConf *kuberniteConfig.Config) error {
	gitRepo, err := git.NewRepositoryFromFilePath(kuberniteConf.DeploymentFileRepositoryPath)
	if err != nil {
		return err
	}
	err = gitRepo.CommitDeployment(kuberniteConf.DeploymentFileRepositoryPath, kuberniteConf.DeploymentFilePath, kuberniteConf.GitUsername, kuberniteConf.GitPassword)
	if err != nil {
		return err
	}
	return nil
}

func handleDeployment(kuberniteConf *kuberniteConfig.Config) (*kubernetesManifest.Deployment, error) {
	switch kuberniteConf.BuildEvent {
	case git.TagEvent:
		return updateDeploymentForTagEvent(kuberniteConf)
	default:
		return updateDeploymentForOtherEvent(kuberniteConf)
	}
}

func updateDeploymentForTagEvent(kuberniteConf *kuberniteConfig.Config) (*kubernetesManifest.Deployment, error) {
	// open git repository
	gitRepo, err := git.NewRepositoryFromFilePath(kuberniteConf.DeploymentTagRepositoryPath)
	if err != nil {
		return nil, err
	}

	// get the latest git tag on the repository
	latestTag, err := gitRepo.GetLatestTagName()
	if err != nil {
		return nil, err
	}

	// open deployment file
	deploymentFile, err := kubernetesManifest.NewDeploymentFromFile(kuberniteConf.DeploymentFilePath)
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

	if err := deploymentFile.UpdateImageTag(kuberniteConf.DeploymentImageName, latestTag); err != nil {
		log.Fatal(err)
	}

	return deploymentFile, nil
}

func updateDeploymentForOtherEvent(kuberniteConf *kuberniteConfig.Config) (*kubernetesManifest.Deployment, error) {
	// open git repository
	gitRepo, err := git.NewRepositoryFromFilePath(kuberniteConf.DeploymentTagRepositoryPath)
	if err != nil {
		log.Fatal(err)
	}

	// get the latest commit hash in the repository
	latestCommitHash, err := gitRepo.GetLatestCommitHash()
	if err != nil {
		return nil, err
	}

	// open deployment file
	deploymentFile, err := kubernetesManifest.NewDeploymentFromFile(kuberniteConf.DeploymentFilePath)
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
	if err := deploymentFile.UpdateImageTag(kuberniteConf.DeploymentImageName, "latest"); err != nil {
		log.Fatal(err)
	}

	return deploymentFile, nil
}

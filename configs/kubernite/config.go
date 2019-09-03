package kubernite

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
	"kubernite/pkg/git"
)

func init() {
	err := viper.BindEnv("KubernetesServer", "PLUGIN_KUBERNETES_SERVER")
	err = viper.BindEnv("KubernetesCertData", "PLUGIN_KUBERNETES_CERT_DATA")
	err = viper.BindEnv("KubernetesClientCertData", "PLUGIN_KUBERNETES_CLIENT_CERT_DATA")
	err = viper.BindEnv("KubernetesClientKeyData", "PLUGIN_KUBERNETES_CLIENT_KEY_DATA")
	err = viper.BindEnv("KubernetesDeploymentFilePath", "PLUGIN_KUBERNETES_DEPLOYMENT_FILE_PATH")
	err = viper.BindEnv("DeploymentTagRepositoryPath", "PLUGIN_DEPLOYMENT_TAG_REPOSITORY_PATH")
	err = viper.BindEnv("DeploymentImageName", "PLUGIN_DEPLOYMENT_IMAGE_NAME")
	err = viper.BindEnv("DryRun", "PLUGIN_DRY_RUN")
	err = viper.BindEnv("DeploymentFileRepositoryPath", "PLUGIN_DEPLOYMENT_FILE_REPOSITORY_PATH")
	err = viper.BindEnv("CommitDeployment", "PLUGIN_COMMIT_DEPLOYMENT")
	err = viper.BindEnv("BuildEvent", "DRONE_BUILD_EVENT")
	if err != nil {
		err = ErrPackageInitialisation{Reasons: []string{
			"binding viper keys to environment variables",
			err.Error(),
		}}
		log.Fatal(err)
	}
}

type Config struct {
	//// TODO insert again these add `validate:"required"` to top 5!!
	KubernetesServer             string
	KubernetesCertData           string
	KubernetesClientCertData     string
	KubernetesClientKeyData      string
	KubernetesDeploymentFilePath string
	DeploymentTagRepositoryPath  string
	DeploymentImageName          string
	DryRun                       bool
	DeploymentFileRepositoryPath string
	CommitDeployment             bool
	BuildEvent                   git.Event `validate:"required"`
}

func GetConfig() (*Config, error) {
	// set default configuration
	//viper.SetDefault("DeploymentTagRepositoryPath", "/drone/src") TODO - add again
	//viper.SetDefault("DryRun", false) TODO - add again
	viper.SetDefault("CommitDeployment", false)

	// TODO remove
	viper.SetDefault("BuildEvent", "push")
	viper.SetDefault("DeploymentTagRepositoryPath", "/Users/Lawrence/go/src/github.com/andile-innovation/james")
	viper.SetDefault("KubernetesDeploymentFilePath", "/Users/Lawrence/go/src/github.com/andile-innovation/konductor/manifests/dev/james/deployment.yaml")
	viper.SetDefault("DryRun", true)
	//viper.SetDefault("DeploymentImageName", "tbcloud/james")

	// parse the config from environment
	conf := new(Config)
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, err
	}

	//// TODO insert again these
	// validate the configuration
	if err := validator.New().Struct(conf); err != nil {
		return nil, ErrInvalidConfig{Reasons: []string{err.Error()}}
	}

	return conf, nil
}

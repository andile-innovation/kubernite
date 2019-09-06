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
	err = viper.BindEnv("DeploymentFilePath", "PLUGIN_DEPLOYMENT_FILE_PATH")
	err = viper.BindEnv("DeploymentTagRepositoryPath", "PLUGIN_DEPLOYMENT_TAG_REPOSITORY_PATH")
	err = viper.BindEnv("DeploymentImageName", "PLUGIN_DEPLOYMENT_IMAGE_NAME")
	err = viper.BindEnv("DryRun", "PLUGIN_DRY_RUN")
	err = viper.BindEnv("DeploymentFileRepositoryPath", "PLUGIN_DEPLOYMENT_FILE_REPOSITORY_PATH")
	err = viper.BindEnv("CommitDeployment", "PLUGIN_COMMIT_DEPLOYMENT")
	err = viper.BindEnv("BuildEvent", "DRONE_BUILD_EVENT")
	//TODO - used for git push
	//err = viper.BindEnv("GitUsername", "PLUGIN_GIT_USERNAME")
	//err = viper.BindEnv("GitPassword", "PLUGIN_GIT_PASSWORD")
	//err = viper.BindEnv("GitKey", "PLUGIN_GIT_KEY")
	if err != nil {
		err = ErrPackageInitialisation{Reasons: []string{
			"binding viper keys to environment variables",
			err.Error(),
		}}
		log.Fatal(err)
	}
}

type Config struct {
	KubernetesServer             string `validate:"required"`
	KubernetesCertData           string `validate:"required"`
	KubernetesClientCertData     string `validate:"required"`
	KubernetesClientKeyData      string `validate:"required"`
	DeploymentFilePath           string `validate:"required"`
	DeploymentTagRepositoryPath  string
	DeploymentImageName          string
	DryRun                       bool
	DeploymentFileRepositoryPath string
	CommitDeployment             bool
	BuildEvent                   git.Event `validate:"required"`
	//TODO - used for git push
	//GitPassword                  string
	//GitUsername                  string
	//GitKey                       string
}

func GetConfig() (*Config, error) {
	// set default configuration
	viper.SetDefault("DeploymentTagRepositoryPath", "/drone/src")
	viper.SetDefault("DryRun", false)

	// parse the config from environment
	conf := new(Config)
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, err
	}

	// validate the configuration
	if err := validator.New().Struct(conf); err != nil {
		return nil, ErrInvalidConfig{Reasons: []string{err.Error()}}
	}

	return conf, nil
}

package kubernite

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
)

func init() {
	err := viper.BindEnv("DeploymentRepositoryPath", "PLUGIN_DEPLOYMENT_REPOSITORY_PATH")
	err = viper.BindEnv("KubernetesDeploymentFilePath", "PLUGIN_KUBERNETES_DEPLOYMENT_FILE_PATH")
	err = viper.BindEnv("KubernetesServer", "PLUGIN_KUBERNETES_SERVER")
	err = viper.BindEnv("KubernetesCertData", "PLUGIN_KUBERNETES_CERT_DATA")
	err = viper.BindEnv("KubernetesClientCertData", "PLUGIN_KUBERNETES_CLIENT_CERT_DATA")
	err = viper.BindEnv("KubernetesClientKeyData", "PLUGIN_KUBERNETES_CLIENT_KEY_DATA")
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
	DeploymentRepositoryPath     string `validate:"required"`
	KubernetesDeploymentFilePath string `validate:"required"`
	KubernetesServer             string `validate:"required"`
	KubernetesCertData           string `validate:"required"`
	KubernetesClientCertData     string `validate:"required"`
	KubernetesClientKeyData      string `validate:"required"`
	BuildEvent                   string `validate:"required,eq=push|eq=tag|eq=merge"`
}

func GetConfig() (*Config, error) {
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

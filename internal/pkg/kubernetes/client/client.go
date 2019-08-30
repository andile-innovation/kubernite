package client

import (
	"k8s.io/client-go/kubernetes"
	kubernetesRestClient "k8s.io/client-go/rest"
	kuberniteConfig "kubernite/configs/kubernite"
)

type Client struct {
	*kubernetes.Clientset
}

func NewClientFromKuberniteConfig(kuberniteConf *kuberniteConfig.Config) (*Client, error) {
	// create rest client configuration from kubernite config
	var restClientConfig = new(kubernetesRestClient.Config)
	restClientConfig.Host = kuberniteConf.KubernetesServer
	restClientConfig.TLSClientConfig.CAData = []byte(kuberniteConf.KubernetesCertData)
	restClientConfig.TLSClientConfig.CertData = []byte(kuberniteConf.KubernetesClientCertData)
	restClientConfig.TLSClientConfig.KeyData = []byte(kuberniteConf.KubernetesClientKeyData)

	// create the kubernetes client set
	clientset, err := kubernetes.NewForConfig(restClientConfig)
	if err != nil {
		return nil, ErrCreatingClientSet{Reasons: []string{
			err.Error(),
		}}
	}

	return &Client{
		Clientset: clientset,
	}, nil
}

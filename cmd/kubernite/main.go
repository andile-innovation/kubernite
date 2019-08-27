package main

import (
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	kubeRestclient "k8s.io/client-go/rest"
	"os"
)

func main() {
	var config = new(kubeRestclient.Config)
	var err error

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

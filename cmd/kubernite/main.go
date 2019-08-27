package main

import (
	"encoding/base64"
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	kubeRestclient "k8s.io/client-go/rest"
	"log"
	"os"
)

func main() {
	var config = new(kubeRestclient.Config)
	var err error

	config.Host = os.Getenv("PLUGIN_KUBERNETES_SERVER")
	config.TLSClientConfig.CAData, err = base64.URLEncoding.DecodeString(os.Getenv("PLUGIN_KUBERNETES_CERT_DATA"))
	if err != nil {
		log.Fatal("error decoding CA Data: " + err.Error())
	}
	config.TLSClientConfig.CertData, err = base64.URLEncoding.DecodeString(os.Getenv("PLUGIN_CLIENT_CERT_DATA"))
	if err != nil {
		log.Fatal("error decoding Cert Data: " + err.Error())
	}
	config.TLSClientConfig.KeyData, err = base64.URLEncoding.DecodeString(os.Getenv("PLUGIN_CLIENT_KEY_DATA"))
	if err != nil {
		log.Fatal("error decoding Key Data: " + err.Error())
	}

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

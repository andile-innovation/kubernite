package main

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"io"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	kubeRestclient "k8s.io/client-go/rest"
	"log"
	"os"
)

var repoPath = "/Users/bernardbussy/go/src/github.com/andile-innovation/james"

func main() {

	var config = new(kubeRestclient.Config)
	var err error

	r, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatal("error opening git repo: " + err.Error())
	}

	deploymentTag := ""

	tagReferences, err := r.Tags()
	if err != nil {
		log.Fatal("error getting repo tag refs: " + err.Error())
	}

	if tagReferences != nil {
		latestTag, err := tagReferences.Next()
		switch err {
		case io.EOF:
		case nil:
			deploymentTag = latestTag.Name().Short()
		default:
			log.Fatal("error getting latest tag ref: " + err.Error())
		}
	}

	if deploymentTag == "" {
		commitReferences, err := r.CommitObjects()
		if err != nil {
			log.Fatal("error getting repo commits: " + err.Error())
		}
		commit, err := commitReferences.Next()
		if err != nil {
			log.Fatal("error getting latest commit: " + err.Error())
		}
		deploymentTag = commit.Hash.String()
	}

	fmt.Println("use: " + deploymentTag)

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

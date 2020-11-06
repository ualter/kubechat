package main

import (
	"flag"
	"fmt"
	"path/filepath"

	//"k8s.io/apimachinery/pkg/api/errors"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"	
	"k8s.io/client-go/util/homedir"
)

var (
	masterURL  string
	kubeconfig string
)

func main() {
	fmt.Printf("Kubeconfig is %s\n",kubeconfig)
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

func init() {
	if home := homedir.HomeDir(); home != "" {
		flag.StringVar(&kubeconfig,"kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		flag.StringVar(&kubeconfig,"kubeconfig", "~/.kube/config", "absolute path to the kubeconfig file")
	}
	//flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}

package main

import (
	"context"
	"flag"
	"path/filepath"
	"time"
	"log"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	cgv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"	
	"k8s.io/client-go/util/homedir"
	corev1 "k8s.io/api/core/v1"
)

var (
	masterURL  string
	kubeconfig string
	api cgv1.CoreV1Interface
)

func main() {
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	api = clientset.CoreV1()

	for {
		listPods(nil)
		options := metav1.ListOptions{LabelSelector:"app=teachstore-course"}
		fmt.Println("")
		listPods(&options)

		/*podFound := getPod("steachstore-course-1.0.0-674b855dc-2j787","develop")
		if podFound != nil {
			log.Printf("podFond=%s",podFound.Name)
		}*/

		time.Sleep(10 * time.Second)
	}
}

func listPods(listOptions *metav1.ListOptions) {
	if listOptions == nil {
	   log.Printf("\033[0;33mListing All PODs\033[0;0m")
	   listOptions = &metav1.ListOptions{}
	} else {
	   log.Printf("\033[0;33mListing PODs with Label \033[0;36m%s\033[0;0m", listOptions.LabelSelector)	
	}
	pods, err := api.Pods("").List(context.TODO(), *listOptions)
	if err != nil {
		panic(err.Error())
	}
	log.Printf("\033[0;33mFound \033[0;36m%d\033[0;0m running\033[0;0m\n",len(pods.Items))
	for _, pod := range pods.Items {
		//log.Println(pod.Spec.Containers[0].Image)
		log.Println(pod.Name)
	}
}

func getPod(podName string, namespace string) *corev1.Pod {
	pod, err := api.Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		log.Printf("Pod %s in namespace %s not found\n", podName, namespace)
		return nil
	} else if err != nil {
		panic(err.Error())
	} else {
		log.Printf("Pod %s in namespace %s found!", podName, namespace)
	}
	return pod
}

func init() {
	if home := homedir.HomeDir(); home != "" {
		flag.StringVar(&kubeconfig,"kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		flag.StringVar(&kubeconfig,"kubeconfig", "~/.kube/config", "absolute path to the kubeconfig file")
	}
}

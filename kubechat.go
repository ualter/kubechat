package main

import (
	"context"
	"flag"
	"path/filepath"
	"log"
	"fmt"
	"time"

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


	listOptions := metav1.ListOptions{LabelSelector:"app=teachstore-course"}
	pods := findPodByOptions(&listOptions,"develop")
	log.Printf("Total \033[0;33m%d\033[0;0m Pods found", len(pods.Items))
	for idx, p := range pods.Items {
		//log.Printf("\033[0;34%d\033[0;0 - \033[0;33 %s\033[0;0",idx,p.Name)
		log.Printf("\033[0;33m[%d]\033[0;0m - \033[0;36m%s\033[0;0m",idx+1,p.Name)
	}
	

	//findPodTeachStoreCourse()
	//listPodsEachSeconds()
}

func findPodTeachStoreCourse() {
	podFoundByName := findPodByName("teachstore-course-1.0.0-674b855dc-2j787","develop")
	if podFoundByName != nil {
		log.Printf("POD %s \033[0;33mFOUND\033[0;0m!",podFoundByName.Name)
	}
}

func listPodsEachSeconds() {
	// List Pods each 10 seconds
	for {
		// All of them
		listPods(nil)
		fmt.Printf("")

		// With LabelSelector
		options := metav1.ListOptions{LabelSelector:"app=teachstore-course"}
		fmt.Println("")
		listPods(&options)

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

func findPodByName(podName string, namespace string) *corev1.Pod {
	pod, err := api.Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})

	if errors.IsNotFound(err) {
		log.Printf("Pod %s in namespace %s \033[0;33mNOT FOUND\033[0;0m!\n", podName, namespace)
		return nil
	} else if err != nil {
		panic(err.Error())
	}/* else {
		log.Printf("Pod %s in namespace %s found!", podName, namespace)
	}*/
	return pod
}

func findPodByOptions(listOptions *metav1.ListOptions, namespace string) *corev1.PodList {
	pods, err := api.Pods(namespace).List(context.TODO(), *listOptions)
	if err != nil {
		panic(err.Error())
	}
	return pods
}

func init() {
	if home := homedir.HomeDir(); home != "" {
		flag.StringVar(&kubeconfig,"kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		flag.StringVar(&kubeconfig,"kubeconfig", "~/.kube/config", "absolute path to the kubeconfig file")
	}
}

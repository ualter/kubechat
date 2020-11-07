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
	kubernetes "k8s.io/client-go/kubernetes"
	cgv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"	
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/homedir"
	kubeinformers "k8s.io/client-go/informers"
	corev1 "k8s.io/api/core/v1"
	//appsv1 "k8s.io/api/apps/v1"
)

var (
	masterURL  string
	kubeconfig string
	api cgv1.CoreV1Interface
)

func main() {
	flag.Parse()

	// kubeconfig := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	api = clientset.CoreV1()

	createPod(clientset,"develop")

	//watchPods(clientset)
	//findPodByOptions	
	//findPodByLabelSelector()
	// findPodTeachStoreCourse()
	//listPodsEachSeconds()
}

func watchPods(clientset *kubernetes.Clientset) {
	factory := kubeinformers.NewSharedInformerFactory(clientset, 0)
	informer := factory.Core().V1().Pods().Informer()
	stopper := make(chan struct{})
	defer close(stopper)
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			mObj := obj.(metav1.Object)
			log.Printf("New Pod Added to Store: %s", mObj.GetName())
		},
		UpdateFunc: func(old, new interface{}) {
			mObj := new.(metav1.Object)
			log.Printf("Update Pod at Store: %s", mObj.GetName())
		},
	})
	informer.Run(stopper)
}

func createPod(clientset *kubernetes.Clientset, namespace string) {
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind: "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-pod",
			Namespace: "develop",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
				    Image: "nginx",
				},
			},
		},
	}
	result, err := api.Pods(namespace).Create(context.TODO(),pod,metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created Pod %s",result.Name)
}

func findPodByLabelSelector() {
	listOptions := metav1.ListOptions{
		LabelSelector:"app=teachstore-course",
	}
	//FieldSelector: fmt.Sprintf("spec.ports[0].nodePort=%s", port),
	pods := findPodByOptions(&listOptions,"develop")
	log.Printf("Total \033[0;33m%d\033[0;0m Pods found", len(pods.Items))
	for idx, p := range pods.Items {
		//log.Printf("\033[0;34%d\033[0;0 - \033[0;33 %s\033[0;0",idx,p.Name)
		log.Printf("\033[0;33m[%d]\033[0;0m - \033[0;36m%s\033[0;0m",idx+1,p.Name)
	}
}

func findPodTeachStoreCourse() {
	podFoundByName := findPodByName("teachstore-course-1.0.0-674b855dc-4l6nr","develop")
	if podFoundByName != nil {
		log.Printf("POD %s \033[0;33mFOUND\033[0;0m!",podFoundByName.Name)
	}
	container  := podFoundByName.Spec.Containers[0]
	fmt.Printf("%s\n",container.Image)
	/*
	containers := podFoundByName.Spec.Containers
	fmt.Printf("%s",&container)
	fmt.Printf("%s\n",container.Name)
	fmt.Printf("%s\n",container.Command)
	*/
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


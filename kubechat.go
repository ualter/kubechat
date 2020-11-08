package main

import (
	"context"
	"flag"
	"path/filepath"
	"log"
	"fmt"
	"time"
	"encoding/json"
	coreerrors "errors"

	jsonpatch "github.com/mattbaird/jsonpatch"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkgtypes "k8s.io/apimachinery/pkg/types"
	kubernetes "k8s.io/client-go/kubernetes"
	cgcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	cgappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/tools/clientcmd"	
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"
	kubeinformers "k8s.io/client-go/informers"
	corev1 "k8s.io/api/core/v1"
	//appsv1 "k8s.io/api/apps/v1"
)

var (
	masterURL  string
	kubeconfig string
	apiCoreV1 cgcorev1.CoreV1Interface
	apiAppsV1 cgappsv1.AppsV1Interface
)

type PatchPathSpec struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

func main() {
	flag.Parse()
	clientset := startK8sClient()
	_ = clientset

	/*pod := findPodByName("teachstore-course-1.0.0-674b855dc-hw9wt","develop")
	if pod != nil {
		log.Printf("POD %s \033[0;33mFOUND\033[0;0m!",pod.Name)
	}
	log.Printf("%s",pod.Name)
	log.Printf("%s",pod.Spec.EphemeralContainers)*/

	
	
	//updateDeployment(clientset)
	//listDeployment(clientset)
	//createPod(clientset,"develop")
	//watchPods(clientset)
	//findPodByOptions	
	//findPodByLabelSelector()
	//findPodTeachStoreCourse()
	//listPodsEachSeconds()
	//applyPatchDeploymentWithReplicas(Int32Ptr(1))
	applyPatchDeploymentAddContainer()
	//applyPatchDeploymentRemoveContainer()
}

func startK8sClient() *kubernetes.Clientset {
	// kubeconfig := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	apiCoreV1 = clientset.CoreV1()
	apiAppsV1 = clientset.AppsV1()
	return clientset
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

func listDeployment(clientset *kubernetes.Clientset) {
	deployments, err := apiAppsV1.Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error());
	}
	log.Printf("Total of Deployments...: \033[0;33m%d\033[0;0m",len(deployments.Items))
	for idx, d := range deployments.Items {
		log.Printf("[%d] \033[0;33m%s\033[0;0m",idx+1,d.Name)
		log.Printf("    Labels:")
		for k, v := range d.Spec.Selector.MatchLabels {
			log.Printf(" \033[0;34m-->\033[0;33m %s\033[0;36m=\033[0;33m%s\033[0;0m",k,v)
		}
	}
}

func updateDeployment(clientset *kubernetes.Clientset) {
	//TODO: those, it will be the entry arguments
	labels    := "app=teachstore-course"
	namespace := "develop"
	replicas  := Int32Ptr(1)

	if len(namespace) == 0 {
		panic(coreerrors.New("The namespace must be informed"))
	}

	options := metav1.ListOptions{
		LabelSelector: labels,
	}
	deployments, err := apiAppsV1.Deployments(namespace).List(context.TODO(), options)
	if err != nil {
		panic(err.Error());
	}
	if len(deployments.Items) > 0 {
		log.Printf("Total of Deployments to Update...: \033[0;33m%d\033[0;0m",len(deployments.Items))
		for idx, deploy := range deployments.Items {
			log.Printf("[%d] \033[0;33m%s\033[0;0m",idx+1,deploy.Name)

			retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
				result, getErr := apiAppsV1.Deployments(namespace).Get(context.TODO(), deploy.Name, metav1.GetOptions{})
				if getErr != nil {
					panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
				}

				result.Spec.Replicas = replicas
			    // result.Spec.Template.Spec.Containers[0].Image = "nginx:1.13" // change nginx version
				_, updateErr := apiAppsV1.Deployments(namespace).Update(context.TODO(), result, metav1.UpdateOptions{})
				return updateErr
			})
			if retryErr != nil {
				log.Panic("Update failed: %v", retryErr)
				panic(retryErr)
			}
			log.Printf("Updated deployment done!")
		}
	} else {
		log.Printf("No deployment found with the label: %s",labels)
	}
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
	result, err := apiCoreV1.Pods(namespace).Create(context.TODO(),pod,metav1.CreateOptions{})
	if err != nil {
		panic(err.Error());
	}
	fmt.Printf("Created Pod %s",result.Name)
}

func findPodByLabelSelector() *corev1.PodList {
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
	return pods
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
	pods, err := apiCoreV1.Pods("").List(context.TODO(), *listOptions)
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
	pod, err := apiCoreV1.Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})

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
	pods, err := apiCoreV1.Pods(namespace).List(context.TODO(), *listOptions)
	if err != nil {
		panic(err.Error())
	}
	return pods
}

func applyPatchDeploymentWithReplicas(numberOfReplicas *int32) {
	// Read the Deployment to be Patched
	deployment, err := apiAppsV1.Deployments("develop").Get(context.TODO(), "teachstore-course-1.0.0",metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	// Converts Actual Deployment to JSON
	jsonDeploymentBefore, err := json.Marshal(deployment)
	if err != nil {
		panic(err.Error())
	}
	// Change the Deployment (modification)
	deployment.Spec.Replicas = numberOfReplicas
	jsonDeploymentAfter, err := json.Marshal(deployment)
	if err != nil {
		panic(err.Error())
	}
	// Create a JSON Patch (http://jsonpatch.com/ JSON Patch is specified in RFC 6902 from the IETF) - using library for that
	patch, err := jsonpatch.CreatePatch(jsonDeploymentBefore, jsonDeploymentAfter)
	if err != nil {
		log.Fatalln(err)
	}
	patchBytes, err := json.MarshalIndent(patch, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(patchBytes))
	/*ko := &appsv1.Deployment{}
	ko.Spec.Template.Spec.Containers = append(ko.Spec.Template.Spec.Containers, corev1.Container{
		Name:            "busybox-new",
		Image:           "busybox",
	})*/
    // Apply the Patch
	result, err := apiAppsV1.Deployments("develop").Patch(context.TODO(),deployment.Name,pkgtypes.JSONPatchType,patchBytes,metav1.PatchOptions{})
	if err != nil {
		panic(err.Error())
	} else {
		log.Printf("%s",result)
	}
}

func applyPatchDeploymentAddContainer() {
	// Read the Deployment to be Patched
	deployment, err := apiAppsV1.Deployments("develop").Get(context.TODO(), "teachstore-course-1.0.0",metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	// Converts Actual Deployment to JSON
	jsonDeploymentBefore, err := json.Marshal(deployment)
	if err != nil {
		panic(err.Error())
	}
	// Change the Deployment (modification) - Adding a New Container
	deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, corev1.Container{
		Name:            "busybox",
		Image:           "busybox",
		ImagePullPolicy: corev1.PullIfNotPresent,
		Command: []string{
			"sleep",
			"3600",
		},
	})
	jsonDeploymentAfter, err := json.Marshal(deployment)
	if err != nil {
		panic(err.Error())
	}
	// Create a JSON Patch (http://jsonpatch.com/ JSON Patch is specified in RFC 6902 from the IETF) - using library for that
	patch, err := jsonpatch.CreatePatch(jsonDeploymentBefore, jsonDeploymentAfter)
	if err != nil {
		log.Fatalln(err)
	}
	patchBytes, err := json.MarshalIndent(patch, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(patchBytes))
	
    // Apply the Patch
	result, err := apiAppsV1.Deployments("develop").Patch(context.TODO(),deployment.Name,pkgtypes.JSONPatchType,patchBytes,metav1.PatchOptions{})
	if err != nil {
		panic(err.Error())
	} else {
		log.Printf("%s",result)
	}
	
}

func applyPatchDeploymentRemoveContainer() {
	// Read the Deployment to be Patched
	deployment, err := apiAppsV1.Deployments("develop").Get(context.TODO(), "teachstore-course-1.0.0",metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	// Converts Actual Deployment to JSON
	jsonDeploymentBefore, err := json.Marshal(deployment)
	if err != nil {
		panic(err.Error())
	}
	// Change the Deployment (modification) - Remove Container
	deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers[0:1])
	jsonDeploymentAfter, err := json.Marshal(deployment)
	if err != nil {
		panic(err.Error())
	}
	// Create a JSON Patch (http://jsonpatch.com/ JSON Patch is specified in RFC 6902 from the IETF) - using library for that
	patch, err := jsonpatch.CreatePatch(jsonDeploymentBefore, jsonDeploymentAfter)
	if err != nil {
		log.Fatalln(err)
	}
	patchBytes, err := json.MarshalIndent(patch, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(patchBytes))
	
    // Apply the Patch
	result, err := apiAppsV1.Deployments("develop").Patch(context.TODO(),deployment.Name,pkgtypes.JSONPatchType,patchBytes,metav1.PatchOptions{})
	if err != nil {
		panic(err.Error())
	} else {
		log.Printf("%s",result)
	}
	
}

func init() {
	if home := homedir.HomeDir(); home != "" {
		flag.StringVar(&kubeconfig,"kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		flag.StringVar(&kubeconfig,"kubeconfig", "~/.kube/config", "absolute path to the kubeconfig file")
	}
}


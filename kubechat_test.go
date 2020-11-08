package main

import (
	"testing"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestFindPodByLabelSelector(t *testing.T) {
	startK8sClient()
	pods := findPodByLabelSelector()
	if pods == nil {
		t.Error("Ops... the List of Pods are null :-(")
	}
}

func TestFindPodByName(t *testing.T) {
	startK8sClient()
	findPodByName("teachstore-course-1.0.0-674b855dc-4l6nr","develop")
}

func TestListPods(t *testing.T) {
	startK8sClient()
	listPods(nil)
}

func TestListPodsByOptions(t *testing.T) {
	startK8sClient()
	options := metav1.ListOptions{LabelSelector:"app=teachstore-course"}
	listPods(&options)
}


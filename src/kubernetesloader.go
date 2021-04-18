package main

import (
	"encoding/json"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func getK8sClient() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}

//func loadPods(clientset kubernetes.Clientset) *apiv1.PodList {
// get pods in all the namespaces by omitting namespace
// Or specify namespace to get pods in particular namespace
//pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
//if err != nil {
//	panic(err.Error())
//}

//return pods

func loadTestdata(name string, obj interface{}) {
	file, err := os.Open("testdata/" + name + ".json")
	if err != nil {
		panic(err.Error())
	}

	dec := json.NewDecoder(file)
	dec.Decode(&obj)
}

func loadPods(ns string) v1.PodList {
	var obj v1.PodList
	loadTestdata("pods", &obj)
	return obj
}

func loadServices(ns string) v1.ServiceList {
	var obj v1.ServiceList
	loadTestdata("services", &obj)
	return obj
}

func loadStatefulSets(ns string) appsv1.StatefulSetList {
	var obj appsv1.StatefulSetList
	loadTestdata("statefulsets", &obj)
	return obj
}

func loadReplicaSets(ns string) appsv1.ReplicaSetList {
	var obj appsv1.ReplicaSetList
	loadTestdata("replicasets", &obj)
	return obj
}

func loadDaemonSets(ns string) appsv1.DaemonSetList {
	var obj appsv1.DaemonSetList
	loadTestdata("daemonsets", &obj)
	return obj
}

func loadIngresses(ns string) extensionsv1beta1.IngressList {
	var obj extensionsv1beta1.IngressList
	loadTestdata("ingresses", &obj)
	return obj
}

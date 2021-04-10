package main

import (
	"context"
	"encoding/json"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func loadPods(clientset kubernetes.Clientset) *apiv1.PodList {
	// get pods in all the namespaces by omitting namespace
	// Or specify namespace to get pods in particular namespace
	pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	return pods
}

func loadServices() apiv1.ServiceList {
	file, err := os.Open("testdata/services.json")
	if err != nil {
		panic(err.Error())
	}

	dec := json.NewDecoder(file)
	var obj apiv1.ServiceList
	dec.Decode(&obj)
	return obj
}

func loadDaemonSets() appsv1.DaemonSetList {
	file, err := os.Open("testdata/daemonsets.json")
	if err != nil {
		panic(err.Error())
	}

	dec := json.NewDecoder(file)
	var obj appsv1.DaemonSetList
	dec.Decode(&obj)
	return obj
}

func loadDeployments() appsv1.DeploymentList {
	file, err := os.Open("testdata/deployments.json")
	if err != nil {
		panic(err.Error())
	}

	dec := json.NewDecoder(file)
	var obj appsv1.DeploymentList
	dec.Decode(&obj)
	return obj
}

func loadStatefulSets() appsv1.StatefulSetList {
	file, err := os.Open("testdata/statefulsets.json")
	if err != nil {
		panic(err.Error())
	}

	dec := json.NewDecoder(file)
	var obj appsv1.StatefulSetList
	dec.Decode(&obj)
	return obj
}

func loadReplicaSets() appsv1.ReplicaSetList {
	file, err := os.Open("testdata/replicasets.json")
	if err != nil {
		panic(err.Error())
	}

	dec := json.NewDecoder(file)
	var obj appsv1.ReplicaSetList
	dec.Decode(&obj)
	return obj
}

func loadIngresses() extensionsv1beta1.IngressList {
	file, err := os.Open("testdata/ingresses.json")
	if err != nil {
		panic(err.Error())
	}

	dec := json.NewDecoder(file)
	var obj extensionsv1beta1.IngressList
	dec.Decode(&obj)
	return obj
}

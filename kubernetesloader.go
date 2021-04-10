package main

import (
	"encoding/json"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
)

func loadPods() apiv1.PodList {
	file, err := os.Open("testdata/pods.json")
	if err != nil {
		panic(err.Error())
	}

	dec := json.NewDecoder(file)
	var obj apiv1.PodList
	dec.Decode(&obj)
	return obj
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

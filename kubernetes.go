package main

import (
	"encoding/json"
	"fmt"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/selection"
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

func subgraph(name string, graph string) string {
	return fmt.Sprintf("subgraph %s\n%send\n", name, graph)
}

func podNode(obj apiv1.Pod) string {
	return fmt.Sprintf("POD%s\n", obj.ObjectMeta.Name)
}

func serviceNode(obj apiv1.Service) string {
	return fmt.Sprintf("SVC%s\n", obj.ObjectMeta.Name)
}

func statefulSetNode(obj appsv1.StatefulSet) string {
	return fmt.Sprintf("SSET%s\n", obj.ObjectMeta.Name)
}

func replicaSetNode(obj appsv1.ReplicaSet) string {
	return fmt.Sprintf("RSET%s\n", obj.ObjectMeta.Name)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsKey(s map[string]string, e string) bool {
	for a, _ := range s {
		if a == e {
			return true
		}
	}
	return false
}

func selectPods(pods []apiv1.Pod, matchExpressions []v1.LabelSelectorRequirement, matchLabels map[string]string) []string {
	var foundPods []string

	for k, v := range matchLabels {
		matchExpressions = append(matchExpressions, v1.LabelSelectorRequirement{
			Key:      k,
			Operator: v1.LabelSelectorOperator(selection.In),
			Values:   []string{v},
		})
	}

	for _, pod := range pods {
		labels := pod.ObjectMeta.Labels
		matched := true

		for _, expression := range matchExpressions {
			key := expression.Key
			values := expression.Values
			op := expression.Operator

			switch op {
			case "in":
				if !containsKey(labels, key) {
					matched = false
					break
				}
				for _, search := range values {
					if labels[key] == search {
						break
					}
					matched = false
					break
				}
			case "NotIn":
				// TODO
			case "Exists":
				if !containsKey(labels, key) {
					matched = false
					break
				}
			case "DoesNotExist":
				if containsKey(labels, key) {
					matched = false
					break
				}
			}
		}

		if matched {
			foundPods = append(foundPods, pod.ObjectMeta.Name)
		}
	}

	return foundPods
}

func labelsSelected(labels map[string]string, selectors map[string]string) bool {
	for key, value := range selectors {
		foundValue, present := labels[key]
		if !present {
			return false
		}
		if foundValue == value {
			continue
		}
		return false
	}
	return true
}

func selectSets(sets map[string]map[string]string, selectors map[string]string) []string {
	if len(selectors) == 0 {
		return make([]string, 0)
	}

	var foundSets []string

	for setName, set := range sets {
		if labelsSelected(set, selectors) {
			foundSets = append(foundSets, setName)
		}
	}

	return foundSets
}
func selectStatefulSets(sets []appsv1.StatefulSet, selectors map[string]string) []string {
	setMap := make(map[string]map[string]string, len(sets))
	for _, set := range sets {
		if set.Status.Replicas > 0 {
			setMap[set.ObjectMeta.Name] = set.ObjectMeta.Labels
		}
	}
	return selectSets(setMap, selectors)
}
func selectReplicaSets(sets []appsv1.ReplicaSet, selectors map[string]string) []string {
	setMap := make(map[string]map[string]string, len(sets))
	for _, set := range sets {
		if set.Status.Replicas > 0 {
			setMap[set.ObjectMeta.Name] = set.ObjectMeta.Labels
		}
	}
	return selectSets(setMap, selectors)
}

func generateChart(podList apiv1.PodList, serviceList apiv1.ServiceList, statefulSetList appsv1.StatefulSetList, replicaSetList appsv1.ReplicaSetList) string {
	chart := "graph LR\n"

	podGraph := ""
	serviceGraph := ""
	setGraph := ""
	edges := ""

	for _, obj := range podList.Items {
		podGraph += podNode(obj)
	}
	for _, obj := range serviceList.Items {
		serviceGraph += serviceNode(obj)
		for _, selectedPod := range selectStatefulSets(statefulSetList.Items, obj.Spec.Selector) {
			edges += fmt.Sprintf("SVC%s --> SSET%s\n", obj.ObjectMeta.Name, selectedPod)
		}
		for _, selectedPod := range selectReplicaSets(replicaSetList.Items, obj.Spec.Selector) {
			edges += fmt.Sprintf("SVC%s --> RSET%s\n", obj.ObjectMeta.Name, selectedPod)
		}
	}
	for _, obj := range statefulSetList.Items {
		setGraph += statefulSetNode(obj)
		for _, selectedPod := range selectPods(podList.Items, obj.Spec.Selector.MatchExpressions, obj.Spec.Selector.MatchLabels) {
			edges += fmt.Sprintf("SSET%s --> POD%s\n", obj.ObjectMeta.Name, selectedPod)
		}
	}
	for _, obj := range replicaSetList.Items {
		setGraph += replicaSetNode(obj)
		for _, selectedPod := range selectPods(podList.Items, obj.Spec.Selector.MatchExpressions, obj.Spec.Selector.MatchLabels) {
			edges += fmt.Sprintf("RSET%s --> POD%s\n", obj.ObjectMeta.Name, selectedPod)
		}
	}

	return chart +
		subgraph("Pods", podGraph) +
		subgraph("Services", serviceGraph) +
		subgraph("Sets", setGraph) +
		edges
}

func main() {
	fmt.Println(generateChart(loadPods(), loadServices(), loadStatefulSets(), loadReplicaSets()))
}

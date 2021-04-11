package main

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/selection"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsKeyStrStr(s map[string]string, e string) bool {
	for a, _ := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsKeyStrInt(s map[string]int, e string) bool {
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
				if !containsKeyStrStr(labels, key) {
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
				if !containsKeyStrStr(labels, key) {
					matched = false
					break
				}
			case "DoesNotExist":
				if containsKeyStrStr(labels, key) {
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

func selectServices(ingressName string, services []apiv1.Service, rules []extensionsv1beta1.IngressRule) []string {
	var response []string
	for _, rule := range rules {
		for _, path := range rule.HTTP.Paths {
			response = append(response, fmt.Sprintf("ING%s -- %s%s --> SVC%s%v\n", ingressName, rule.Host, path.Path, path.Backend.ServiceName, path.Backend.ServicePort.IntVal))
		}
	}
	return response
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

func selectStatefulSets(sets []appsv1.StatefulSet, selectors map[string]string) []appsv1.StatefulSet {
	setLabelMap := make(map[string]map[string]string, len(sets))
	setMap := make(map[string]appsv1.StatefulSet, len(sets))
	for _, set := range sets {
		if set.Status.Replicas > 0 {
			setLabelMap[set.ObjectMeta.Name] = set.ObjectMeta.Labels
			setMap[set.ObjectMeta.Name] = set
		}
	}
	var response = make([]appsv1.StatefulSet, 0)
	for _, selected := range selectSets(setLabelMap, selectors) {
		response = append(response, setMap[selected])
	}
	return response
}

func selectReplicaSets(sets []appsv1.ReplicaSet, selectors map[string]string) []appsv1.ReplicaSet {
	setLabelMap := make(map[string]map[string]string, len(sets))
	setMap := make(map[string]appsv1.ReplicaSet, len(sets))
	for _, set := range sets {
		if set.Status.Replicas > 0 {
			setLabelMap[set.ObjectMeta.Name] = set.ObjectMeta.Labels
			setMap[set.ObjectMeta.Name] = set
		}
	}
	var response = make([]appsv1.ReplicaSet, 0)
	for _, selected := range selectSets(setLabelMap, selectors) {
		response = append(response, setMap[selected])
	}
	return response
}

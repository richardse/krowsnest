package main

import (
	"fmt"
	"html"
	"sort"
	"strings"

	v1 "k8s.io/api/core/v1"
)

func trimServiceName(svcName string, ns string) string {
	return strings.ReplaceAll(svcName, fmt.Sprintf("-%s", ns), "")
}

func fS(str string) string {
	// Replaces bad Mermaid strings
	str = strings.ReplaceAll(str, "--", "DH")
	str = strings.ReplaceAll(str, "-x", "x")
	return strings.ReplaceAll(str, "default", "defau1t")
}

func subgraph(name string, graph map[string]bool) string {
	graphSlice := make([]string, len(graph))
	i := 0
	for k, _ := range graph {
		graphSlice[i] = k
		i++
	}
	sort.Strings(graphSlice)

	graphstr := ""
	for _, g := range graphSlice {
		graphstr += g
	}
	return fmt.Sprintf("subgraph %s\n%send\n", name, graphstr)
}

func selectorToString(selectors map[string]string) string {
	responses := make([]string, len(selectors))
	i := 0
	for k, v := range selectors {
		if len(k) > 18 && k[0:18] == "app.kubernetes.io/" {
			k = k[18:]
		}
		responses[i] = fmt.Sprintf("%s=%s", k, v)
		i++
	}
	sort.Strings(responses)
	return strings.Join(responses, ",<br/>")
}

func mapEventsToNodeNames(events v1.EventList) map[string]string {
	mappedEvents := make(map[string][]string)
	for _, event := range events.Items {
		nodeName := fmt.Sprintf("%s%s", event.InvolvedObject.Kind, event.InvolvedObject.Name)
		eventTime := ""
		if !event.EventTime.IsZero() {
			eventTime = event.EventTime.String()
		} else if !event.LastTimestamp.IsZero() {
			eventTime = event.LastTimestamp.String()
		} else {
			eventTime = event.FirstTimestamp.String()
		}
		eventStr := html.EscapeString(fmt.Sprintf("%s - %s: %s", eventTime, event.Reason, event.Message))
		mappedEvents[nodeName] = append(mappedEvents[nodeName], eventStr)
	}

	mappedStringEvents := make(map[string]string, len(mappedEvents))
	for k, v := range mappedEvents {
		mappedStringEvents[k] = fmt.Sprintf(" <span style='color:deepskyblue;' title='%s'>[?]</span>", strings.Join(v, "\n"))
	}

	return mappedStringEvents
}

func generateChart(objects KubernetesObjects) string {
	chart := "graph LR\n"

	podGraph := make(map[string]bool)
	serviceGraph := make(map[string]bool)
	setGraph := make(map[string]bool)
	ingressGraph := make(map[string]bool)
	edges := make([]string, 0)

	eventsByNode := mapEventsToNodeNames(objects.events)

	for _, obj := range objects.pods.Items {
		podGraph[podNode(obj, eventsByNode[fmt.Sprintf("Pod%s", obj.ObjectMeta.Name)])] = true
		if obj.OwnerReferences != nil && len(obj.OwnerReferences) > 0 && obj.OwnerReferences[0].Kind == "DaemonSet" {
			edges = append(edges, fmt.Sprintf("DSET%s --> POD%s", obj.OwnerReferences[0].Name, obj.ObjectMeta.Name))
		}
	}

	for _, obj := range objects.daemonsets.Items {
		setGraph[daemonSetNode(obj, eventsByNode[fmt.Sprintf("DaemonSet%s", obj.ObjectMeta.Name)])] = true
	}

	loadBalanced := false
	for _, obj := range objects.services.Items {
		for _, port := range obj.Spec.Ports {
			serviceGraph[serviceNode(obj, port, eventsByNode[fmt.Sprintf("Service%s", obj.ObjectMeta.Name)])] = true

			if obj.Spec.Type == "LoadBalancer" {
				loadBalanced = true
				edges = append(edges, fmt.Sprintf("INGCloudLoadBalancer --> SVC%s%v", obj.ObjectMeta.Name, port.Port))
			}

			found := false
			for _, selectedSet := range selectStatefulSets(objects.statefulsets.Items, obj.Spec.Selector) {
				setGraph[statefulSetNode(selectedSet, eventsByNode[fmt.Sprintf("StatefulSet%s", selectedSet.ObjectMeta.Name)])] = true
				edges = append(edges, fmt.Sprintf("SVC%s%d -- %s --> SSET%s", obj.ObjectMeta.Name, port.Port, selectorToString(obj.Spec.Selector), selectedSet.ObjectMeta.Name))
				found = true
			}
			for _, selectedSet := range selectReplicaSets(objects.replicasets.Items, obj.Spec.Selector) {
				setGraph[replicaSetNode(selectedSet, eventsByNode[fmt.Sprintf("ReplicaSet%s", selectedSet.ObjectMeta.Name)])] = true
				edges = append(edges, fmt.Sprintf("SVC%s%d -- %s --> RSET%s", obj.ObjectMeta.Name, port.Port, selectorToString(obj.Spec.Selector), selectedSet.ObjectMeta.Name))
				found = true
			}

			if !found {
				for _, selectedPod := range selectPods(objects.pods.Items, nil, obj.Spec.Selector) {
					edges = append(edges, fmt.Sprintf("SVC%s%d -- %s --> POD%s", obj.ObjectMeta.Name, port.Port, selectorToString(obj.Spec.Selector), selectedPod))
				}
			}
		}
	}

	if loadBalanced {
		ingressGraph["INGCloudLoadBalancer(\"Cloud LoadBalancer\")\n"] = true
	}

	for _, obj := range objects.statefulsets.Items {
		for _, selectedPod := range selectPods(objects.pods.Items, obj.Spec.Selector.MatchExpressions, obj.Spec.Selector.MatchLabels) {
			edges = append(edges, fmt.Sprintf("SSET%s --> POD%s", obj.ObjectMeta.Name, selectedPod))
		}
	}
	for _, obj := range objects.replicasets.Items {
		for _, selectedPod := range selectPods(objects.pods.Items, obj.Spec.Selector.MatchExpressions, obj.Spec.Selector.MatchLabels) {
			edges = append(edges, fmt.Sprintf("RSET%s --> POD%s", obj.ObjectMeta.Name, selectedPod))
		}
	}
	for _, obj := range objects.ingresses.Items {
		ingressGraph[ingressNode(obj, eventsByNode[fmt.Sprintf("Ingress%s", obj.ObjectMeta.Name)])] = true
		edges = append(edges, selectServices(obj.ObjectMeta.Name, objects.services.Items, obj.Spec.Rules)...)
	}

	sort.Strings(edges)

	return chart +
		subgraph("Ingresses", ingressGraph) +
		subgraph("Pods", podGraph) +
		subgraph("Services", serviceGraph) +
		subgraph("Sets", setGraph) +
		strings.Join(edges, "\n")
}

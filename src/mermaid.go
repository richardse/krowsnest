package main

import (
	"fmt"
	"html"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
)

var stateToColourMap = map[string]string{"Running": "green", "Ready": "green", "Completed": "blue"}

func stateToColour(state string) string {
	colour, pres := stateToColourMap[state]
	if pres {
		return colour
	}
	return "red"
}

func trimServiceName(svcName string, ns string) string {
	return strings.ReplaceAll(svcName, fmt.Sprintf("-%s", ns), "")
}

func fS(str string) string {
	// Replaces bad Mermaid strings
	str = strings.ReplaceAll(str, "--", "DH")
	str = strings.ReplaceAll(str, "-x", "x")
	return strings.ReplaceAll(str, "default", "defau1t")
}

func subgraph(name string, graph string) string {
	return fmt.Sprintf("subgraph %s\n%send\n", name, graph)
}

func podNode(obj apiv1.Pod) string {
	statusListCount := make(map[string]int)
	statusListMessage := make(map[string]string)

	if obj.Status.ContainerStatuses != nil {
		for _, c := range obj.Status.ContainerStatuses {
			if c.Ready {
				statusListCount["Running"] += 1
				statusListMessage["Running"] = "Running"
			} else if c.State.Terminated != nil {
				reason := "Terminated"
				if c.State.Terminated.Reason != "" {
					reason = c.State.Terminated.Reason
				}
				statusListCount[reason] += 1
				if c.State.Terminated.Message != "" {
					statusListMessage[reason] = c.State.Terminated.Message
				} else {
					statusListMessage[reason] = c.State.Terminated.Reason
				}
			} else if c.State.Waiting != nil {
				reason := "Waiting"
				if c.State.Waiting.Reason != "" {
					reason = c.State.Waiting.Reason
				}
				statusListCount[reason] += 1
				if c.State.Waiting.Message != "" {
					statusListMessage[reason] = c.State.Waiting.Message
				} else {
					statusListMessage[reason] = c.State.Waiting.Reason
				}
			} else if c.State.Running != nil {
				statusListCount["NotReady"] += 1
				statusListMessage["NotReady"] = "NotReady"
			}
		}

		if len(statusListCount) == 1 && containsKeyStrInt(statusListCount, "ContainerCreating") {
			for _, c := range obj.Status.Conditions {
				if c.Type == "ContainersReady" {
					statusListMessage["ContainerCreating"] = fmt.Sprintf("%s: %s", c.Reason, c.Message)
				}
			}
		}
	} else {
		// There are no running containers
		if obj.Status.Conditions != nil && len(obj.Status.Conditions) > 0 {
			condition := obj.Status.Conditions[0]
			statusListCount[condition.Reason] = 1
			statusListMessage[condition.Reason] = condition.Message
		} else {
			statusListCount["NoContainers"] = 1
			statusListMessage["NoContainers"] = "No containers"
		}
	}

	total := 0
	for _, i := range statusListCount {
		total += i
	}

	message := make([]string, len(statusListCount))
	for state, count := range statusListCount {
		message = append(message, fmt.Sprintf("<span style='color:%s' title='%s'>%s %v/%v</span>", stateToColour(state), html.EscapeString(statusListMessage[state]), html.EscapeString(state), count, total))
	}

	return podOrJobNode(obj, strings.Join(message, "<br />"))
}

func podOrJobNode(obj apiv1.Pod, content string) string {
	name := fS(obj.ObjectMeta.Name)

	return fmt.Sprintf("POD%s(\"%s%s\")\n", name, name, content)
}

func serviceNode(obj apiv1.Service, port apiv1.ServicePort) string {
	svcName := trimServiceName(obj.ObjectMeta.Name, obj.ObjectMeta.Namespace)
	portName := ""
	if port.Name != "" {
		portName = fmt.Sprintf(" (%s)", port.Name)
	} else {
		portName = ""
	}

	return fmt.Sprintf("SVC%s%d(\"%s:%d%s\")\n", obj.ObjectMeta.Name, port.Port, svcName, port.Port, portName)
}

func statefulSetNode(obj appsv1.StatefulSet) string {
	return fmt.Sprintf("SSET%s(\"%s\")\n", obj.ObjectMeta.Name, obj.ObjectMeta.Name)
}

func replicaSetNode(obj appsv1.ReplicaSet) string {
	replicas := fmt.Sprintf("<br />Replicas: %d/%d", obj.Status.Replicas, *obj.Spec.Replicas)
	if *obj.Spec.Replicas > obj.Status.Replicas {
		message := "Creation Pending"
		if obj.Status.Conditions != nil {
			for _, set := range obj.Status.Conditions {
				if set.Status == "True" {
					message = strings.ReplaceAll(set.Message, "\"", " ")
				}
			}
		}
		replicas = fmt.Sprintf("<br /><a style='color: red' title='%s'>%s {?}</a>", message, replicas)
	}

	return fmt.Sprintf("RSET%s(\"%s%s\")\n", obj.ObjectMeta.Name, obj.ObjectMeta.Name, replicas)
}

func ingressNode(obj extensionsv1beta1.Ingress) string {
	return fmt.Sprintf("ING%s(\"%s\")\n", obj.ObjectMeta.Name, obj.ObjectMeta.Name)
}

func containerNode(pod apiv1.Pod, container apiv1.Container) string {
	for _, c := range pod.Status.ContainerStatuses {
		if c.Name == container.Name {
			var statusMessage string
			if c.State.Running != nil {
				statusMessage = fmt.Sprintf("Running since %s", c.State.Running.StartedAt)
			} else if c.State.Terminated != nil {
				statusMessage = fmt.Sprintf("Terminated since %s", c.State.Running.StartedAt)
			} else if c.State.Waiting != nil {
				statusMessage = c.State.Waiting.Message
			}
			return fmt.Sprintf("CNT%s%s(\"%s<br />%s\")\n", pod.ObjectMeta.Name, container.Name, container.Name, statusMessage)
		}
	}

  return ""
}

func selectorToString(selectors map[string]string) string {
	responses := make([]string, len(selectors))
	i := 0
	for k, v := range selectors {
		if k[0:18] == "app.kubernetes.io/" {
			k = k[18:]
		}
		responses[i] = fmt.Sprintf("%s=%s", k, v)
		i++
	}
	return strings.Join(responses, ",<br/>")
}

func generateChart(podList apiv1.PodList, serviceList apiv1.ServiceList, statefulSetList appsv1.StatefulSetList, replicaSetList appsv1.ReplicaSetList, ingressList extensionsv1beta1.IngressList) string {

	chart := "graph LR\n"

	podGraph := ""
	serviceGraph := ""
	setGraph := ""
	ingressGraph := ""
	edges := ""

	for _, obj := range podList.Items {
		podGraph += podNode(obj)
	}

	loadBalanced := false
	for _, obj := range serviceList.Items {
		for _, port := range obj.Spec.Ports {
			serviceGraph += serviceNode(obj, port)

			if obj.Spec.Type == "LoadBalancer" {
				loadBalanced = true
				edges += fmt.Sprintf("INGCloudLoadBalancer --> SVC%s%v\n", obj.ObjectMeta.Name, port.Port)
			}

			found := false
			for _, selectedSet := range selectStatefulSets(statefulSetList.Items, obj.Spec.Selector) {
				setGraph += statefulSetNode(selectedSet)
				edges += fmt.Sprintf("SVC%s%d -- %s --> SSET%s\n", obj.ObjectMeta.Name, port.Port, selectorToString(obj.Spec.Selector), selectedSet.ObjectMeta.Name)
				found = true
			}
			for _, selectedSet := range selectReplicaSets(replicaSetList.Items, obj.Spec.Selector) {
				setGraph += replicaSetNode(selectedSet)
				edges += fmt.Sprintf("SVC%s%d -- %s --> RSET%s\n", obj.ObjectMeta.Name, port.Port, selectorToString(obj.Spec.Selector), selectedSet.ObjectMeta.Name)
				found = true
			}

			if !found {
				for _, selectedPod := range selectPods(podList.Items, nil, obj.Spec.Selector) {
					edges += fmt.Sprintf("SVC%s%d -- %s --> POD%s\n", obj.ObjectMeta.Name, port.Port, selectorToString(obj.Spec.Selector), selectedPod)
				}
			}
		}
	}
	if loadBalanced {
		ingressGraph += "INGCloudLoadBalancer(\"Cloud LoadBalancer\")\n"
	}

	for _, obj := range statefulSetList.Items {
		for _, selectedPod := range selectPods(podList.Items, obj.Spec.Selector.MatchExpressions, obj.Spec.Selector.MatchLabels) {
			edges += fmt.Sprintf("SSET%s --> POD%s\n", obj.ObjectMeta.Name, selectedPod)
		}
	}
	for _, obj := range replicaSetList.Items {
		for _, selectedPod := range selectPods(podList.Items, obj.Spec.Selector.MatchExpressions, obj.Spec.Selector.MatchLabels) {
			edges += fmt.Sprintf("RSET%s --> POD%s\n", obj.ObjectMeta.Name, selectedPod)
		}
	}
	for _, obj := range ingressList.Items {
		ingressGraph += ingressNode(obj)
		for _, ingressToServiceEdge := range selectServices(obj.ObjectMeta.Name, serviceList.Items, obj.Spec.Rules) {
			edges += ingressToServiceEdge
		}
	}

	return chart +
		subgraph("Ingresses", ingressGraph) +
		subgraph("Pods", podGraph) +
		subgraph("Services", serviceGraph) +
		subgraph("Sets", setGraph) +
		edges
}

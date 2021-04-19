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

func podNode(obj apiv1.Pod, events string) string {
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

	return podOrJobNode(obj, strings.Join(message, "<br />"), events)
}

func podOrJobNode(obj apiv1.Pod, content string, events string) string {
	name := fS(obj.ObjectMeta.Name)

	return fmt.Sprintf("POD%s(\"%s%s%s\")\n", name, name, events, content)
}

func serviceNode(obj apiv1.Service, port apiv1.ServicePort, events string) string {
	svcName := trimServiceName(obj.ObjectMeta.Name, obj.ObjectMeta.Namespace)
	portName := ""
	if port.Name != "" {
		portName = fmt.Sprintf(" (%s)", port.Name)
	} else {
		portName = ""
	}

	return fmt.Sprintf("SVC%s%d(\"%s:%d%s%s\")\n", obj.ObjectMeta.Name, port.Port, svcName, port.Port, portName, events)
}

func statefulSetNode(obj appsv1.StatefulSet, events string) string {
	replicas := fmt.Sprintf("Replicas: %d/%d", obj.Status.Replicas, *obj.Spec.Replicas)
	if *obj.Spec.Replicas > obj.Status.Replicas {
		message := "Creation Pending"
		if obj.Status.Conditions != nil {
			for _, set := range obj.Status.Conditions {
				if set.Status == "True" {
					message = html.EscapeString(set.Message)
				}
			}
		}
		replicas = fmt.Sprintf("<a style='color: red' title='%s'>%s {?}</a>", message, replicas)
	}

	return fmt.Sprintf("SSET%s(\"%s<br />%s%s\")\n", obj.ObjectMeta.Name, obj.ObjectMeta.Name, replicas, events)
}

func replicaSetNode(obj appsv1.ReplicaSet, events string) string {
	replicas := fmt.Sprintf("Replicas: %d/%d", obj.Status.Replicas, *obj.Spec.Replicas)
	if *obj.Spec.Replicas > obj.Status.Replicas {
		message := "Creation Pending"
		if obj.Status.Conditions != nil {
			for _, set := range obj.Status.Conditions {
				if set.Status == "True" {
					message = html.EscapeString(set.Message)
				}
			}
		}
		replicas = fmt.Sprintf("<a style='color: red' title='%s'>%s {?}</a>", message, replicas)
	}

	return fmt.Sprintf("RSET%s(\"%s<br />%s%s\")\n", obj.ObjectMeta.Name, obj.ObjectMeta.Name, replicas, events)
}

func daemonSetNode(obj appsv1.DaemonSet, events string) string {
	replicas := fmt.Sprintf("Replicas: %d/%d", obj.Status.NumberReady, obj.Status.DesiredNumberScheduled)
	if obj.Status.DesiredNumberScheduled > obj.Status.NumberReady {
		message := "Creation Pending"
		if obj.Status.Conditions != nil {
			for _, set := range obj.Status.Conditions {
				if set.Status == "True" {
					message = html.EscapeString(set.Message)
				}
			}
		}
		replicas = fmt.Sprintf("<a style='color: red' title='%s'>%s {?}</a>", message, replicas)
	}

	return fmt.Sprintf("DSET%s(\"%s<br />%s%s\")\n", obj.ObjectMeta.Name, obj.ObjectMeta.Name, replicas, events)
}

func ingressNode(obj extensionsv1beta1.Ingress, events string) string {
	return fmt.Sprintf("ING%s(\"%s%s\")\n", obj.ObjectMeta.Name, obj.ObjectMeta.Name, events)
}

func containerNode(pod apiv1.Pod, container apiv1.Container, events string) string {
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
			return fmt.Sprintf("CNT%s%s(\"%s<br />%s%s\")\n", pod.ObjectMeta.Name, container.Name, container.Name, statusMessage, events)
		}
	}

	return ""
}

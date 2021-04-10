package main

import (
	"testing"
)

func TestMermaid(t *testing.T) {
	mermaid := generateChart(loadPods(), loadServices(), loadStatefulSets(), loadReplicaSets())
	if mermaid != `graph LR
subgraph Pods
PODsvclb-my-release-wordpress-7dqrn
PODmy-release-mariadb-0
PODmy-release-wordpress-67855bb4bc-rzrfd
end
subgraph Services
SVCkubernetes
SVCmy-release-mariadb
SVCmy-release-wordpress
end
subgraph Sets
SSETmy-release-mariadb
RSETmy-release-wordpress-67855bb4bc
end
SVCmy-release-mariadb --> SSETmy-release-mariadb
SVCmy-release-wordpress --> RSETmy-release-wordpress-67855bb4bc
SSETmy-release-mariadb --> PODmy-release-mariadb-0
RSETmy-release-wordpress-67855bb4bc --> PODmy-release-wordpress-67855bb4bc-rzrfd
` {
		t.Errorf("Mermaid chart does not match expected value.")
	}
}

func TestServiceString(t *testing.T) {
  mermaid = ""
  portName := portName != "" ? $" ({portName})" : ""
  if mermaid != `SVC{mermaidName}{port}(\"<a href='https://dashboard/#/service/{ns}/{name}?namespace={ns}'>{name}:{port}{portName}</a><br /><a style='float: right' href='{krowsnesturl}/api/config/{testcluster}/services/{name}/name'><img src='images/json.jpg'></img></a>` {
    t.Errorf("Service String does not match expected value.")
  }
}
func TestServiceString(t *testing.T) {
}
func TestServiceString(t *testing.T) {
}
func TestServiceString(t *testing.T) {
}
func TestServiceString(t *testing.T) {
}
func TestServiceString(t *testing.T) {
}

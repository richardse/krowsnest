package main

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

func compareNewlinedStrings(expected string, actual string) bool {
	e := strings.Split(expected, "\n")
	a := strings.Split(actual, "\n")
	fail := false

	if len(e) != len(a) {
		fmt.Printf("Expected length=%d Actual length=%d\n", len(e), len(a))
		fail = true
	}
	for i := 0; i < int(math.Min(float64(len(e)), float64(len(a)))); i++ {
		if e[i] != a[i] {
			fmt.Printf("L%d Exp: %s\nL%d Act: %s\n", i, e[i], i, a[i])
			fail = true
		}
	}
	return fail
}

func TestMermaid(t *testing.T) {
	mermaid := generateChart(loadAllTestData())

	if compareNewlinedStrings(`graph LR
subgraph Ingresses
INGCloudLoadBalancer("Cloud LoadBalancer")
INGexample-ingress("example-ingress")
end
subgraph Pods
PODkrowsnest-756fc9d459-qrqw7("krowsnest-756fc9d459-qrqw7<br /><span style='color:green' title='Running'>Running 1/1</span>")
PODkubernetes("kubernetes<br /><span style='color:green' title='Running'>Running 1/1</span>")
PODmy-release-mariadb-0("my-release-mariadb-0<br /><span style='color:red' title='back-off 2m40s restarting failed container=mariadb pod=my-release-mariadb-0_default(5f36eb40-61a2-4a77-95dc-637a21db3371)'>CrashLoopBackOff 1/1</span>")
PODmy-release-wordpress-67855bb4bc-rzrfd("my-release-wordpress-67855bb4bc-rzrfd<br /><span style='color:red' title='back-off 1m20s restarting failed container=wordpress pod=my-release-wordpress-67855bb4bc-rzrfd_default(1844e083-5db8-4f8f-9e98-f78d6c6e42d7)'>CrashLoopBackOff 1/1</span>")
PODsvclb-my-release-wordpress-7dqrn("svclb-my-release-wordpress-7dqrn<br /><span style='color:red' title='0/1 nodes are available: 1 node(s) didn&#39;t have free ports for the requested pod ports.'>Unschedulable 1/1</span>")
end
subgraph Services
SVCkrowsnest443("krowsnest:443 (https)")
SVCkubernetes443("kubernetes:443 (https)")
SVCmy-release-mariadb3306("my-release-mariadb:3306 (mysql)")
SVCmy-release-wordpress443("my-release-wordpress:443 (https)")
SVCmy-release-wordpress80("my-release-wordpress:80 (http)")
end
subgraph Sets
DSETsvclb-my-release-wordpress("svclb-my-release-wordpress<br /><a style='color: red' title='Creation Pending'>Replicas: 0/1 {?}</a>")
RSETkrowsnest-756fc9d459("krowsnest-756fc9d459<br />Replicas: 1/1")
RSETmy-release-wordpress-67855bb4bc("my-release-wordpress-67855bb4bc<br />Replicas: 1/1")
SSETmy-release-mariadb("my-release-mariadb<br />Replicas: 1/1")
end
DSETsvclb-my-release-wordpress --> PODsvclb-my-release-wordpress-7dqrn
INGCloudLoadBalancer --> SVCmy-release-wordpress443
INGCloudLoadBalancer --> SVCmy-release-wordpress80
INGexample-ingress -- www.example.com/ --> SVCmy-release-wordpress80
INGexample-ingress -- www.example.com/banana --> SVCmy-release-wordpress80
RSETkrowsnest-756fc9d459 --> PODkrowsnest-756fc9d459-qrqw7
RSETmy-release-wordpress-67855bb4bc --> PODmy-release-wordpress-67855bb4bc-rzrfd
SSETmy-release-mariadb --> PODmy-release-mariadb-0
SVCkrowsnest443 -- app=krowsnest,<br/>pod-template-hash=756fc9d459 --> RSETkrowsnest-756fc9d459
SVCkubernetes443 -- instance=pod,<br/>name=kubernetes --> PODkubernetes
SVCmy-release-mariadb3306 -- component=primary,<br/>instance=my-release,<br/>name=mariadb --> SSETmy-release-mariadb
SVCmy-release-wordpress443 -- instance=my-release,<br/>name=wordpress --> RSETmy-release-wordpress-67855bb4bc
SVCmy-release-wordpress80 -- instance=my-release,<br/>name=wordpress --> RSETmy-release-wordpress-67855bb4bc`, mermaid) {
		t.Errorf("Mermaid chart does not match expected value.")
	}
}

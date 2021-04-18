package debug

import (
	"fmt"
)

func main() {
	ns := "default"
	fmt.Println(generateChart(loadPods(ns), loadServices(ns), loadStatefulSets(ns), loadReplicaSets(ns), loadDaemonSets(ns), loadIngresses(ns)))
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"regexp"

	pretty "github.com/kr/pretty"
)

var validStaticPath = regexp.MustCompile("^/static/([a-zA-Z0-9]+).(html|js)$")
var validK8sPath = regexp.MustCompile("^/k8s/(daemonsets|deployments|pods|replicasets|services|statefulsets)$")

func loadK8sData(objtype string) (string, error) {
	filename := "testdata/" + objtype + ".json"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index.html", http.StatusFound)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	m := validStaticPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		fmt.Printf("404 %s\n", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	filename := fmt.Sprintf("wwwroot/%s.%s", m[1], m[2])
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("404 %s\n", filename)
		http.NotFound(w, r)
		return
	}

	fmt.Printf("200 %s\n", filename)
	fmt.Fprintf(w, string(body))
}

func k8sHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	m := validK8sPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}

	body, err := loadK8sData(m[1])
	if err != nil {
		fmt.Println(err.Error())
		http.NotFound(w, r)
		return
	}

	pretty.Fprintf(w, string(body))
}

func mermaidHandler(w http.ResponseWriter, r *http.Request) {
	chart := generateChart(loadPods(), loadServices(), loadStatefulSets(), loadReplicaSets(), loadIngresses())
	fmt.Fprint(w, chart)
}

func main() {
	http.HandleFunc("/index.html", staticHandler)
	http.HandleFunc("/k8s/", k8sHandler)
	http.HandleFunc("/static/", staticHandler)
	http.HandleFunc("/mermaid", mermaidHandler)
	//http.HandleFunc("/", redirectHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

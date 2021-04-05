package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"os"
	"regexp"

	pretty "github.com/kr/pretty"

	apiv1 "k8s.io/api/core/v1"
)

var validTemplatePath = regexp.MustCompile("^/([a-zA-Z0-9]+).html$")
var validK8sPath = regexp.MustCompile("^/k8s/(daemonsets|deployments|pods|replicasets|services|statefulsets)$")

func loadK8sData(objtype string) (string, error) {
	filename := "testdata/" + objtype + ".json"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func templateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	if r.URL.Path == "/" {
		r.URL.Path = "/index.html"
	}
	m := validTemplatePath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}

	filename := "templates/" + m[1]
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

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
	file, err := os.Open("testdata/pods.json")
	if err != nil {
		panic(err.Error())
	}
	dec := json.NewDecoder(file)

	// Parse it into the internal k8s structs
	var dep apiv1.PodList

	dec.Decode(&dep)

	// Dump the struct in case you want to see what it looks like
	pretty.Println(dep.Items[0])
}

func main() {
	http.HandleFunc("/", templateHandler)
	http.HandleFunc("/index.html", templateHandler)
	http.HandleFunc("/k8s/", k8sHandler)
	http.HandleFunc("/mermaid", mermaidHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

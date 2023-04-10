package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func init() {
	http.HandleFunc("/state", stateHandler)
	http.HandleFunc("/world", worldHandler)
	http.HandleFunc("/viewer", viewerHandler)
	printErr(http.ListenAndServe(":8080", nil))
}

func urlParam(r *http.Request, param string) string {
	query := r.URL.Query()
	if query.Has(param) {
		return query.Get(param)
	}

	return ""
}

func created(w http.ResponseWriter, obj any) {
	w.WriteHeader(http.StatusCreated)
	objBytes, err := json.Marshal(obj)
	printErr(err)
	_, err = w.Write(objBytes)
	printErr(err)
}

func badRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}

func conflict(w http.ResponseWriter) {
	w.WriteHeader(http.StatusConflict)
}

func printErr(err error) {
	if err != nil {
		fmt.Println("3", err)
	}
}

var githubName = "sapphire-ai-dev"

func worldPath() string {
	base, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	split := strings.Split(base, "/")
	for i, seg := range split {
		if seg == githubName {
			return strings.Join(split[:i+1], "/") + "/sapphire-core/world"
		}
	}

	return ""
}

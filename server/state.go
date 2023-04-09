package server

import "net/http"

func stateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getStateHandler(w, r)
	}
}

func getStateHandler(w http.ResponseWriter, r *http.Request) {

}

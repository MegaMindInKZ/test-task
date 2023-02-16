package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
)

const stringLength = 12

var Router = mux.NewRouter()

func main() {
	Router.HandleFunc("/generate-salt", generateSalt).Methods("POST")

	Server := &http.Server{
		Handler: Router,
		Addr:    "127.0.0.1:8001",
	}
	Server.ListenAndServe()
}

func generateSalt(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	output := struct {
		Salt string `json:"salt"`
	}{
		Salt: randomString(),
	}

	jsonResp, _ := json.Marshal(output)
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonResp)
}

func randomString() string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, stringLength)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

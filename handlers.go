package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func sayName(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	name := r.PathValue("name")
	w.Write([]byte(fmt.Sprintf("Hello %s!\n", name)))
}

func greetPerson(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("greetings!\n"))
}

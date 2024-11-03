package main

import (
	//"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Handler struct {
}

func get() {

}

func post(buf []byte) {

}

func delete(buf []byte) {

}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)

	buf, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(buf))

}

func main() {
	var h Handler
	s := &http.Server{
		Addr:    ":8080",
		Handler: &h,
	}
	log.Fatal(s.ListenAndServe())

}

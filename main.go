package main

import (
	//"encoding/json"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Handler struct {
}

type Users struct {
    Users []User `json:"users"`
}

var users Users

type User struct {
	Username string `json:"username"`
}

func  Get() string {
	//TODO: возвращать ответ
	jsonResponse, err := json.Marshal(users)
	if err != nil {
		return ""
	}
	return string(jsonResponse)
}

func Post(body []byte) {
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		fmt.Println(err)
		return
	}
	users.Users = append(users.Users, user)
	//TODO: write data in file
}

func Delete(body []byte) {
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		fmt.Println(err)
		return
	}
	var index int
	for i, inMemoryUser := range users.Users {
		if inMemoryUser.Username == user.Username {index = i}
	}
	//TODO: if index  == nil {}
	users.Users = append(users.Users[:index], users.Users[index+1:]...)
	//TODO: write data in file
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch request_method := r.Method; request_method{
		case "GET":
			//TODO: возвращать ответ
			Get()
		case "POST":
			Post(body)
		case "DELETE":
			Delete(body)
	}
	fmt.Print(r.Method)
	

}

func main() {
	var jsonFile, err = os.Open("users.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.Unmarshal(byteValue, &users)
	var h Handler
	s := &http.Server{
		Addr:    ":8080",
		Handler: &h,
	}
	log.Fatal(s.ListenAndServe())

}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

const (
	usersFile = "users.json"
)

type Handler struct {
	mu sync.Mutex // для безопасной работы с файлом
}

type Users struct {
	Users []User `json:"users"`
}

var users Users

type User struct {
	Username string `json:"username"`
}

// writeToFile сохраняет данные в файл
func (h *Handler) writeToFile() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	data, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("error marshaling users: %w", err)
	}

	err = os.WriteFile(usersFile, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

func Get() (string, error) {
	jsonResponse, err := json.Marshal(users)
	if err != nil {
		return "", fmt.Errorf("error marshaling response: %w", err)
	}
	return string(jsonResponse), nil
}

func (h *Handler) Post(body []byte) error {
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return fmt.Errorf("error unmarshaling user: %w", err)
	}

	// Проверка на пустое имя пользователя
	if user.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	// Проверка на дубликаты
	for _, existingUser := range users.Users {
		if existingUser.Username == user.Username {
			return fmt.Errorf("user %s already exists", user.Username)
		}
	}

	users.Users = append(users.Users, user)
	return h.writeToFile()
}

func (h *Handler) Delete(body []byte) error {
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return fmt.Errorf("error unmarshaling user: %w", err)
	}

	found := false
	var index int
	for i, inMemoryUser := range users.Users {
		if inMemoryUser.Username == user.Username {
			index = i
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("user %s not found", user.Username)
	}

	users.Users = append(users.Users[:index], users.Users[index+1:]...)
	return h.writeToFile()
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading body: %v", err), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		response, err := Get()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, response)

	case http.MethodPost:
		if err := h.Post(body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodDelete:
		if err := h.Delete(body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	// Создаем файл, если он не существует
	if _, err := os.Stat(usersFile); os.IsNotExist(err) {
		users = Users{Users: make([]User, 0)}
		data, _ := json.Marshal(users)
		if err := os.WriteFile(usersFile, data, 0644); err != nil {
			log.Fatalf("Error creating users file: %v", err)
		}
	} else {
		// Читаем существующий файл
		jsonFile, err := os.Open(usersFile)
		if err != nil {
			log.Fatalf("Error opening users file: %v", err)
		}
		defer jsonFile.Close()

		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			log.Fatalf("Error reading users file: %v", err)
		}

		if err := json.Unmarshal(byteValue, &users); err != nil {
			log.Fatalf("Error parsing users file: %v", err)
		}
	}

	var h Handler
	s := &http.Server{
		Addr:    ":8080",
		Handler: &h,
	}

	log.Printf("Server starting on %s", s.Addr)
	log.Fatal(s.ListenAndServe())
}
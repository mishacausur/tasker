package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

var processingError = "Ошибка обработки запроса, попробуйте попытку позднее"
var invalidIDError = "Неверный идентификатор задачи"
var internalError = "Внутренняя ошибка сервера"
var alreadyExist = "Задача с таким ID уже существует"

func getTasks(w http.ResponseWriter, r *http.Request) {

	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, processingError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)

	if err != nil {
		fmt.Println("Ошибка записи данных getTasks (61 строка)")
	}
}

func getTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	task, ok := tasks[taskID]

	if !ok {
		http.Error(w, invalidIDError, http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(task)

	if err != nil {
		http.Error(w, internalError, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(response)

	if err != nil {
		fmt.Println("Ошибка записи данных getTask (88 строка)")
	}
}

func postTasks(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buffer bytes.Buffer

	_, err := buffer.ReadFrom(r.Body)

	if err != nil {
		http.Error(w, processingError, http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buffer.Bytes(), &task); err != nil {
		http.Error(w, processingError, http.StatusBadRequest)
		return
	}

	_, ok := tasks[task.ID]

	if ok {
		http.Error(w, alreadyExist, http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	_, ok := tasks[taskID]

	if !ok {
		http.Error(w, invalidIDError, http.StatusBadRequest)
		return
	}

	delete(tasks, taskID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {

	r := chi.NewRouter()

	r.Get("/tasks", getTasks)
	r.Get("/tasks/{id}", getTask)
	r.Post("/tasks", postTasks)
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}

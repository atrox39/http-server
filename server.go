package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Task struct {
	Title       string
	Description string
}

var tasks []Task

// Estructura estandar de un mensaje
type RequestMessage struct {
	Name    string
	Message string
}

// Para mandar la respuesta como JSON
func JSON(message interface{}) []byte {
	msg, err := json.Marshal(message)
	if err != nil {
		return []byte("")
	}
	return msg
}

// Para obtener el cuerpo de la respuesta en JSON
func JSONBody(r *http.Request) interface{} {
	var data interface{}
	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&data)
		if err != nil {
			return nil
		}
		return data
	}
	return nil
}

// Para obtener el cuerpo de la respuesta en JSON
func JSONTask(r *http.Request, task *Task) error {
	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&task)
		if err != nil {
			return err
		}
	}
	return nil
}

// Controladores
func message(w http.ResponseWriter, message RequestMessage, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(JSON(message))
}

// Mensaje Method Not Allowed
func notAllowed(w http.ResponseWriter) {
	message(w, RequestMessage{
		Name:    "Method",
		Message: "Method not allowed",
	}, http.StatusMethodNotAllowed)
}

func notFound(w http.ResponseWriter) {
	message(w, RequestMessage{
		Name:    "NotFound",
		Message: "Page or item not found",
	}, http.StatusNotFound)
}

func badRequest(w http.ResponseWriter) {
	message(w, RequestMessage{
		Name:    "BadRequest",
		Message: "Body bad request",
	}, http.StatusNotFound)
}

func uknownError(w http.ResponseWriter) {
	message(w, RequestMessage{
		Name:    "UknownError",
		Message: "Uknown server error",
	}, http.StatusInternalServerError)
}

func removeTaskIndex(s []Task, index int) []Task {
	return append(s[:index], s[index+1:]...)
}

func main() {
	http.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" { // Si es tipo GET
			res, err := json.Marshal(tasks)
			if err != nil {
				uknownError(w)
			} else {
				w.Write(res)
			}
		} else if r.Method == "POST" { // Si es tipo POST
			var tmp Task
			err := JSONTask(r, &tmp)
			if err != nil {
				badRequest(w)
			} else {
				tasks = append(tasks, tmp)
				w.Write(JSON(tmp))
			}
		} else if r.Method == "DELETE" { // Si es tipo POST
			var body = JSONBody(r).(map[string]interface{})
			if body["id"] != nil {
				id, err := strconv.Atoi(body["id"].(string))
				if err != nil {
					badRequest(w)
				} else {
					tasks = removeTaskIndex(tasks, id-1)
				}
			}
			if body != nil {
				w.Write(JSON(body))
			} else {
				badRequest(w)
			}
		} else { // Si no entonces manda una respuesta
			notAllowed(w)
		}
	})
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {

	})
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

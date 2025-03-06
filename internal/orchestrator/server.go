package orchestrator

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var calculator = NewDistributedCalculator()

// NewRouter создает новый маршрутизатор и регистрирует обработчики маршрутов.
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(recoveryMiddleware)
	r.HandleFunc("/api/v1/calculate", calculateHandler).Methods("POST")
	r.HandleFunc("/api/v1/expressions", getExpressionsHandler).Methods("GET")
	r.HandleFunc("/api/v1/expressions/{id}", getExpressionByIDHandler).Methods("GET")
	r.HandleFunc("/internal/task", getTaskHandler).Methods("GET")
	r.HandleFunc("/internal/task", postTaskResultHandler).Methods("POST")
	r.HandleFunc("/internal/tasks", getTasksHandler).Methods("GET")
	return r
}

// recoveryMiddleware перехватывает все паники и возвращает статус 500.
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// calculateHandler обрабатывает запрос на добавление вычисления арифметического выражения.
func calculateHandler(w http.ResponseWriter, r *http.Request) {
	var req CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	res, _ := calculator.Calculate(req)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		panic(err)
	}
}

// getExpressionsHandler обрабатывает запрос на получение списка выражений.
func getExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	res, _ := calculator.GetExpressions()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		panic(err)
	}
}

// getExpressionByIDHandler обрабатывает запрос на получение выражения по его идентификатору.
func getExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	res, err := calculator.GetExpressionByID(id)
	if err == ErrNotFound {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		panic(err)
	}
}

// getTaskHandler обрабатывает запрос на получение задачи для выполнения.
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	res, err := calculator.GetTask()
	if err == ErrNotFound {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		panic(err)
	}
}

// postTaskResultHandler обрабатывает запрос на прием результата обработки данных.
func postTaskResultHandler(w http.ResponseWriter, r *http.Request) {
	var req TaskResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	err := calculator.PostTaskResult(req)
	if err == ErrNotFound {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// getTasksHandler обрабатывает запрос на получение всех задач (для демонстрации работы)
func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	res, _ := calculator.GetTasks()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		panic(err)
	}
}

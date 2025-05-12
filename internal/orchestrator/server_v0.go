// Package orchestrator содержит реализацию сервера для распределенного вычислителя.
package orchestrator

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// calculateHandler обрабатывает запрос на добавление вычисления арифметического выражения.
func calculateHandlerV0(w http.ResponseWriter, r *http.Request) {
	var req CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	res, _ := calculator.Calculate(req)
	_, err := db.CreateExpressionWithId("", res.ID, CalculateRequest{Expression: req.Expression})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		panic(err)
	}
}

// getExpressionsHandler обрабатывает запрос на получение списка выражений.
func getExpressionsHandlerV0(w http.ResponseWriter, _ *http.Request) {
	res, _ := calculator.GetExpressions()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		panic(err)
	}
}

// getExpressionByIDHandler обрабатывает запрос на получение выражения по его идентификатору.
func getExpressionByIDHandlerV0(w http.ResponseWriter, r *http.Request) {
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
func getTaskHandlerV0(w http.ResponseWriter, _ *http.Request) {
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
func postTaskResultHandlerV0(w http.ResponseWriter, r *http.Request) {
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
func getTasksHandlerV0(w http.ResponseWriter, _ *http.Request) {
	res, _ := calculator.GetTasks()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		panic(err)
	}
}

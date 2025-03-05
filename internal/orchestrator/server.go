package orchestrator

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// generateUUID генерирует UUID версии 4
func generateUUID4() (string, error) {
	u := make([]byte, 16)
	_, err := rand.Read(u)
	if err != nil {
		return "", fmt.Errorf("ошибка при генерации UUID: %v", err)
	}
	// Устанавливаем версии и варианта UUID
	u[6] = (u[6] & 0x0f) | 0x40 // Версия 4
	u[8] = (u[8] & 0x3f) | 0x80 // Вариант RFC 4122
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:]), nil
}

type ExpressionStatus string

const (
	StatusPending   ExpressionStatus = "pending"
	StatusCompleted ExpressionStatus = "completed"
)

type Expression struct {
	Id     string           `json:"id"`
	Status ExpressionStatus `json:"status"`
	Result float64          `json:"result"`
	Tasks  []Task           `json:"tasks"`
}

type Task struct {
	Id            string    `json:"id"`
	Arg1          float64   `json:"arg1"`
	Arg2          float64   `json:"arg2"`
	Operation     string    `json:"operation"`
	OperationTime time.Time `json:"operation_time"`
}

var expressions = make(map[string]*Expression)
var tasks = make(map[string]*Task)
var mu sync.Mutex

// StartCalCulationRequest - структура запроса на создание вычисления
type StartCalCulationRequest struct {
	Expression string `json:"expression"`
}

// StartCalCulationResponse - структура ответа на создание вычисления
type StartCalCulationResponse struct {
	Id string `json:"id"`
}

// StartCalculationHandler обрабатывает запрос на создание вычисления
func StartCalculationHandler(w http.ResponseWriter, req *http.Request) {
	var r StartCalCulationRequest
	err := json.NewDecoder(req.Body).Decode(&r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	id, err := generateUUID4()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	expression := &Expression{
		Id:     id,
		Status: StatusPending,
	}
	mu.Lock()
	expressions[id] = expression
	mu.Unlock()

	resp := StartCalCulationResponse{Id: id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetExpressionsListResponseItem - элемент списка выражений
type GetExpressionsListResponseItem struct {
	Id     string  `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}

// GetExpressionsListResponse - ответ на запрос списка выражений
type GetExpressionsListResponse struct {
	Expressions []GetExpressionsListResponseItem `json:"expressions"`
}

// GetExpressionsListHandler - обработчик запроса списка выражений
func GetExpressionsListHandler(w http.ResponseWriter, req *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var response GetExpressionsListResponse
	for _, expr := range expressions {
		response.Expressions = append(response.Expressions, GetExpressionsListResponseItem{
			Id:     expr.Id,
			Status: string(expr.Status),
			Result: expr.Result,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetExpressionHandlerResponse - ответ на запрос выражения
type GetExpressionHandlerResponse struct {
	Expression GetExpressionsListResponseItem `json:"expression"`
}

// GetExpressionHandler - обработчик запроса выражения
func GetExpressionHandler(w http.ResponseWriter, req *http.Request) {
	parts_path := strings.Split(req.URL.Path, "/")
	id := parts_path[len(parts_path)-1]

	mu.Lock()
	expr, exists := expressions[id]
	mu.Unlock()

	if !exists {
		http.Error(w, "expression not found", http.StatusNotFound)
		return
	}

	resp := GetExpressionHandlerResponse{
		Expression: GetExpressionsListResponseItem{
			Id:     expr.Id,
			Status: string(expr.Status),
			Result: expr.Result,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// GetTaskHandlerResponseItem - элемент ответа на запрос задачи
type GetTaskHandlerResponseItem struct {
	Id            string    `json:"id"`
	Arg1          string    `json:"arg1"`
	Arg2          string    `json:"arg2"`
	Operation     string    `json:"operation"`
	OperationTime time.Time `json:"operation_time"`
}

// GetTaskHandlerResponse - ответ на запрос задачи
type GetTaskHandlerResponse struct {
	Task GetTaskHandlerResponseItem `json:"task"`
}

// GetTaskHandler - обработчик запроса задачи
func GetTaskHandler(w http.ResponseWriter, req *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	for _, task := range tasks {
		resp := GetTaskHandlerResponse{Task: *task}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}

	http.Error(w, "no tasks available", http.StatusNotFound)
}

/*

    Прием результата обработки данных.

    curl --location 'localhost/internal/task' \
    --header 'Content-Type: application/json' \
    --data '{
      "id": 1,
      "result": 2.5
    }'

Коды ответа:

    200 - успешно записан результат
    404 - нет такой задачи
    422 - невалидные данные
    500 - что-то пошло не так
*/

func SaveTaskResultHandler(w http.ResponseWriter, req *http.Request) {
	var r struct {
		Id     string  `json:"id"`
		Result float64 `json:"result"`
	}
	err := json.NewDecoder(req.Body).Decode(&r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	task, exists := tasks[r.Id]
	if !exists {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	expr, exists := expressions[task.Id]
	if !exists {
		http.Error(w, "expression not found", http.StatusNotFound)
		return
	}

	// Update expression result and status
	expr.Result += r.Result
	expr.Status = StatusCompleted

	delete(tasks, r.Id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

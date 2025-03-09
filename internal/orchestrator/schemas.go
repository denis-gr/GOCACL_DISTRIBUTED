// Package orchestrator содержит схемы данных для пакета orchestrator.
package orchestrator

// CalculateRequest Структура для запроса на добавление вычисления арифметического выражения
type CalculateRequest struct {
	Expression string `json:"expression"`
}

// CalculateResponse Структура для ответа на добавление вычисления арифметического выражения
type CalculateResponse struct {
	ID string `json:"id"`
}

// Expression Структура для выражения
type Expression struct {
	ID     string  `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}

// ExpressionsResponse Структура для ответа на получение списка выражений
type ExpressionsResponse struct {
	Expressions []Expression `json:"expressions"`
}

// ExpressionResponse Структура для ответа на получение выражения по его идентификатору
type ExpressionResponse struct {
	Expression Expression `json:"expression"`
}

// Task Структура для задачи
type Task struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int64   `json:"operation_time"`
}

// TaskResponse Структура для ответа на получение задачи для выполнения
type TaskResponse struct {
	Task Task `json:"task"`
}

// TaskResultRequest Структура для запроса на прием результата обработки данных
type TaskResultRequest struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}

// TaskFull Структура для задачи
type TaskFull struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int64   `json:"operation_time"`
	IsBusy        bool    `json:"is_busy"`
}

// TaskFullResponse Структура для ответа на получение задачи для выполнения
type TaskFullResponse struct {
	TasksFull []TaskFull `json:"tasks"`
}

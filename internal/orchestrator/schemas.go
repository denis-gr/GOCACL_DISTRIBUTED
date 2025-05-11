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

// ExpressionDB Структура для выражения в базе данных
type ExpressionDB struct {
	ID         string  `json:"id"`
	Expression string  `json:"expression"`
	Status     string  `json:"status"`
	Result     float64 `json:"result"`
	CreatorId  string  `json:"creator_id"`
}

// UserDB Структура для пользователя в базе данных
type UserDB struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}

// UserCreateForm Структура для создания пользователя
type UserCreateForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserPublic Структура для публичного представления пользователя
type UserPublic struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// UserLoginForm Структура для входа пользователя
type UserLoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Package orchestrator содержит логику для распределенного вычислителя.
package orchestrator

import (
	"errors"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/denis-gr/GOCACL_DISTRIBUTED/pkg/calc"

	"github.com/google/uuid"
)

// ErrNotFound используется для обозначения ошибки, когда элемент не найден.
var ErrNotFound = errors.New("")

// DistributedCalculator представляет распределенный вычислитель.
type DistributedCalculator struct {
	expressions map[string]Expression
	tasks       map[string]Task
	taskBusy    map[string]bool
	resultChans map[string]chan float64
	mu          sync.Mutex
	db          *DB
}

// NewDistributedCalculator создает новый экземпляр DistributedCalculator.
func NewDistributedCalculator(db *DB) *DistributedCalculator {
	return &DistributedCalculator{
		expressions: make(map[string]Expression),
		tasks:       make(map[string]Task),
		taskBusy:    make(map[string]bool),
		resultChans: make(map[string]chan float64),
		db:          db,
	}
}

func (f *DistributedCalculator) createNewTask(a, b float64, ops string) float64 {
	id, _ := uuid.NewV7()
	idStr := id.String()

	var operationTime string
	switch ops {
	case "+":
		operationTime = os.Getenv("TIME_ADDITION_MS")
	case "-":
		operationTime = os.Getenv("TIME_SUBTRACTION_MS")
	case "*":
		operationTime = os.Getenv("TIME_MULTIPLICATIONS_MS")
	case "/":
		operationTime = os.Getenv("TIME_DIVISIONS_MS")
	}
	operationTimeInt, _ := strconv.ParseInt(operationTime, 10, 64)

	resultChan := make(chan float64)
	f.mu.Lock()
	f.tasks[idStr] = Task{
		ID:            idStr,
		Arg1:          a,
		Arg2:          b,
		Operation:     ops,
		OperationTime: operationTimeInt,
	}
	f.resultChans[idStr] = resultChan
	f.mu.Unlock()
	result := <-resultChan
	return result
}

func (f *DistributedCalculator) saveResult(exprID string, res float64, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	expr, exists := f.expressions[exprID]

	if !exists {
		return
	}
	if err != nil {
		expr.Status = err.Error()
	} else {
		expr.Status = "ok"
		expr.Result = res
	}
	f.expressions[exprID] = expr
	err = db.SetResultExpression(exprID, expr.Status, res)
	if err != nil {
		log.Println(err)
	}
}

func (f *DistributedCalculator) calculate(id, expression string) (CalculateResponse, error) {
	f.mu.Lock()
	f.expressions[id] = Expression{
		ID:     id,
		Status: "running",
		Result: 0,
	}
	f.mu.Unlock()

	ops := calc.Operations{
		PlusFunc: func(a, b float64) float64 {
			return calculator.createNewTask(a, b, "+")
		},
		MinusFunc: func(a, b float64) float64 {
			return calculator.createNewTask(a, b, "-")
		},
		MultiplyFunc: func(a, b float64) float64 {
			return calculator.createNewTask(a, b, "*")
		},
		DivideFunc: func(a, b float64) float64 {
			if b == 0 {
				panic("деление на ноль")
			}
			return calculator.createNewTask(a, b, "/")
		},
	}

	resultChan := make(chan struct {
		res float64
		err error
	})

	go func() {
		ans, err := ops.Calc(expression)
		resultChan <- struct {
			res float64
			err error
		}{res: ans, err: err}
	}()

	go func() {
		result := <-resultChan
		f.saveResult(id, result.res, result.err)
	}()

	return CalculateResponse{ID: id}, nil
}

// Calculate выполняет логику для обработки запроса на добавление вычисления арифметического выражения.
func (f *DistributedCalculator) Calculate(req CalculateRequest) (CalculateResponse, error) {
	id, _ := uuid.NewV7()
	idStr := id.String()
	return f.calculate(idStr, req.Expression)
}

// LoadFromDB загружает данные из базы данных.
func (f *DistributedCalculator) LoadFromDB() {
	expressions, err := f.db.GetAllExpressions()
	if err != nil {
		panic(err)
	}
	for _, expr := range expressions {
		f.calculate(expr.ID, expr.Expression)
	}
}

// GetExpressions выполняет логику для обработки запроса на получение списка выражений.
func (f *DistributedCalculator) GetExpressions() (ExpressionsResponse, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	expressions := []Expression{}
	for _, v := range f.expressions {
		expressions = append(expressions, v)
	}
	return ExpressionsResponse{Expressions: expressions}, nil
}

// GetExpressionByID выполняет логику для обработки запроса на получение выражения по его идентификатору.
func (f *DistributedCalculator) GetExpressionByID(id string) (ExpressionResponse, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	expr, ok := f.expressions[id]
	if !ok {
		return ExpressionResponse{}, ErrNotFound
	}
	return ExpressionResponse{
		Expression: expr,
	}, nil
}

// GetTask выполняет логику для обработки запроса на получение задачи для выполнения.
func (f *DistributedCalculator) GetTask() (TaskResponse, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for id, task := range f.tasks {
		_, ok := f.taskBusy[id]
		if !ok {
			f.taskBusy[id] = true
			return TaskResponse{Task: task}, nil
		}
	}
	return TaskResponse{}, ErrNotFound
}

// GetTasks выполняет логику для обработки запроса на получение всех задач.
func (f *DistributedCalculator) GetTasks() (TaskFullResponse, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	tasks := []TaskFull{}
	for id, task := range f.tasks {
		_, isBusy := f.taskBusy[id]
		tasks = append(tasks, TaskFull{
			ID:            task.ID,
			Arg1:          task.Arg1,
			Arg2:          task.Arg2,
			Operation:     task.Operation,
			OperationTime: task.OperationTime,
			IsBusy:        isBusy,
		})
	}
	return TaskFullResponse{TasksFull: tasks}, nil
}

// PostTaskResult выполняет логику для обработки запроса на прием результата обработки данных.
func (f *DistributedCalculator) PostTaskResult(req TaskResultRequest) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	c, ok := f.resultChans[req.ID]
	if !ok {
		return ErrNotFound
	}
	c <- req.Result
	return nil
}

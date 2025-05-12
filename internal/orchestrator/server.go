// Package orchestrator содержит реализацию сервера для распределенного вычислителя.
package orchestrator

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"

	pb "github.com/denis-gr/GOCACL_DISTRIBUTED/internal/gen"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var db, _ = NewDB("db.sqlite3")
var calculator = NewDistributedCalculator(db)

func StartServer(httpAddr string, grpcAddr string) error {
	// Запуск gRPC-сервера на отдельном порту
	grpcServer := grpc.NewServer()
	pb.RegisterOrchestratorServiceServer(grpcServer, &OrchestratorGRPCServer{})
	reflection.Register(grpcServer)

	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", grpcAddr, err)
	}

	go func() {
		log.Printf("Starting gRPC server at %s", grpcAddr)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// Запуск HTTP-сервера на отдельном порту
	router := NewRouter()
	log.Println("Starting HTTP server at", httpAddr)
	err = http.ListenAndServe(httpAddr, router)
	return err
}

func NewRouter() http.Handler {
	grpcServer := grpc.NewServer()
	pb.RegisterOrchestratorServiceServer(grpcServer, &OrchestratorGRPCServer{})
	reflection.Register(grpcServer)

	router := mux.NewRouter()
	router.Use(recoveryMiddleware)
	router.HandleFunc("/api/v1/calculate", calculateHandler).Methods("POST")
	router.HandleFunc("/api/v1/expressions", getExpressionsHandler).Methods("GET")
	router.HandleFunc("/api/v1/expressions/{id}", getExpressionByIDHandler).Methods("GET")
	router.HandleFunc("/internal/task", getTaskHandler).Methods("GET")
	router.HandleFunc("/internal/task", postTaskResultHandler).Methods("POST")
	router.HandleFunc("/internal/tasks", getTasksHandler).Methods("GET")
	router.HandleFunc("/api/v1/register", registerUserHandler).Methods("POST")
	router.HandleFunc("/api/v1/login", loginUserHandler).Methods("POST")

	return router
}

// registerUserHandler обрабатывает запрос на регистрацию нового пользователя.
func registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var req UserCreateForm
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	if _, err := db.CreateUser(req); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// loginUserHandler обрабатывает запрос на вход пользователя.
func loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var req UserLoginForm
	var flag bool
	var err error
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	flag, err = db.CheckUserPassword(req.Username, req.Password)
	if !flag || (err != nil) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	id, err := db.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      id,
		"username": req.Username,
		"exp":      jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})
	tokenString, err := jwt.SignedString([]byte("secret"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
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

func checkJWTToken(r *http.Request) (string, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return "", http.ErrNoCookie
	}
	tokenString = tokenString[len("Bearer "):]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrNoCookie
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["exp"] != nil {
			exp := claims["exp"].(jwt.NumericDate)
			if time.Now().After(exp.Time) {
				return "", http.ErrNoCookie
			}
		}
		if claims["sub"] != nil {
			id := claims["sub"].(string)
			return id, nil
		}
	}
	return "", http.ErrNoCookie
}

// calculateHandler обрабатывает запрос на добавление вычисления арифметического выражения.
func calculateHandler(w http.ResponseWriter, r *http.Request) {
	var req CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	res, _ := calculator.Calculate(req)
	user_id, _ := checkJWTToken(r)
	_, err := db.CreateExpression(user_id, CalculateRequest{Expression: req.Expression})
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
func getExpressionsHandler(w http.ResponseWriter, _ *http.Request) {
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

type OrchestratorGRPCServer struct {
	pb.UnimplementedOrchestratorServiceServer
}

func (s *OrchestratorGRPCServer) GetTask(ctx context.Context, in *pb.Empty) (*pb.TaskResponse, error) {
	task, err := calculator.GetTask()
	if err == ErrNotFound {
		return nil, status.Errorf(codes.NotFound, "task not found") // Добавлено сообщение об ошибке
	}
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to get task: %v", err) // Добавлено сообщение об ошибке
	}
	return &pb.TaskResponse{
		Task: &pb.Task{
			Id:            task.Task.ID,
			Operation:     task.Task.Operation,
			Arg1:          task.Task.Arg1,
			Arg2:          task.Task.Arg2,
			OperationTime: task.Task.OperationTime,
		},
	}, nil
}

func (s *OrchestratorGRPCServer) SendResult(ctx context.Context, in *pb.TaskResultRequest) (*pb.Empty, error) {
	err := calculator.PostTaskResult(TaskResultRequest{
		ID:     in.Id,
		Result: in.Result,
	})
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

// getTaskHandler обрабатывает запрос на получение задачи для выполнения.
func getTaskHandler(w http.ResponseWriter, _ *http.Request) {
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
func getTasksHandler(w http.ResponseWriter, _ *http.Request) {
	res, _ := calculator.GetTasks()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		panic(err)
	}
}

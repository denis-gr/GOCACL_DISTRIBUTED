package orchestrator

import (
	"testing"
)

func TestCreateUser_Success(t *testing.T) {
	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	form := UserCreateForm{
		Username: "testuser",
		Password: "password123",
	}

	user, err := db.CreateUser(form)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.Username != form.Username {
		t.Errorf("expected username %s, got %s", form.Username, user.Username)
	}
	if user.ID == "" {
		t.Error("expected non-empty user ID")
	}
}

func TestCreateUser_UsernameAlreadyExists(t *testing.T) {

	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	form := UserCreateForm{
		Username: "testuser",
		Password: "password123",
	}

	_, err = db.CreateUser(form)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = db.CreateUser(form)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "username already exists" {
		t.Errorf("expected error 'username already exists', got %v", err)
	}
}

func TestCreateUser_InvalidPassword(t *testing.T) {

	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	form := UserCreateForm{
		Username: "testuser",
		Password: "",
	}

	_, err = db.CreateUser(form)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCheckUserPassword(t *testing.T) {

	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	form := UserCreateForm{
		Username: "testuser",
		Password: "password123",
	}

	_, err = db.CreateUser(form)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	isValid, err := db.CheckUserPassword(form.Username, form.Password)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !isValid {
		t.Error("expected password to be valid, got invalid")
	}

	isValid, err = db.CheckUserPassword(form.Username, "wrongpassword")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if isValid {
		t.Error("expected password to be invalid, got valid")
	}
}

func TestGetUserByUsername(t *testing.T) {
	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	form := UserCreateForm{
		Username: "testuser",
		Password: "password123",
	}

	user, err := db.CreateUser(form)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	retrievedUser, err := db.GetUserByUsername(user.Username)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if retrievedUser.ID != user.ID {
		t.Errorf("expected ID %s, got %s", user.ID, retrievedUser.ID)
	}
	if retrievedUser.Username != user.Username {
		t.Errorf("expected username %s, got %s", user.Username, retrievedUser.Username)
	}
}

func TestGetUserAll(t *testing.T) {

	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	users := []UserCreateForm{
		{Username: "user1", Password: "password1"},
		{Username: "user2", Password: "password2"},
	}
	for _, form := range users {
		_, err := db.CreateUser(form)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	allUsers, err := db.GetUserAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(allUsers) != len(users) {
		t.Errorf("expected %d users, got %d", len(users), len(allUsers))
	}

	for i, user := range allUsers {
		if user.Username != users[i].Username {
			t.Errorf("expected username %s, got %s", users[i].Username, user.Username)
		}
		if user.ID == "" {
			t.Error("expected non-empty user ID")
		}
	}

	userIDs := make(map[string]bool)
	for _, user := range allUsers {
		if _, exists := userIDs[user.ID]; exists {
			t.Errorf("duplicate user ID found: %s", user.ID)
		}
		userIDs[user.ID] = true
	}
}

func TestCreateExpression(t *testing.T) {

	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	// Создаем пользователя для получения creatorID
	userForm := UserCreateForm{
		Username: "testuser",
		Password: "password123",
	}
	user, err := db.CreateUser(userForm)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	form := CalculateRequest{
		Expression: "2+2",
	}

	expression, err := db.CreateExpression(user.ID, form)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if expression.Expression != form.Expression {
		t.Errorf("expected expression %s, got %s", form.Expression, expression.Expression)
	}
	if expression.Status != "running" {
		t.Errorf("expected status 'running', got %s", expression.Status)
	}
	if expression.ID == "" {
		t.Error("expected non-empty expression ID")
	}
}

func TestGetExpressionByID(t *testing.T) {

	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	// Создаем пользователя для получения creatorID
	userForm := UserCreateForm{
		Username: "testuser",
		Password: "password123",
	}
	user, err := db.CreateUser(userForm)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	form := CalculateRequest{
		Expression: "2+2",
	}
	expression, err := db.CreateExpression(user.ID, form)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	retrievedExpression, err := db.GetExpressionByID(expression.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if retrievedExpression.ID != expression.ID {
		t.Errorf("expected ID %s, got %s", expression.ID, retrievedExpression.ID)
	}
	if retrievedExpression.Expression != expression.Expression {
		t.Errorf("expected expression %s, got %s", expression.Expression, retrievedExpression.Expression)
	}
	if retrievedExpression.Status != expression.Status {
		t.Errorf("expected status %s, got %s", expression.Status, retrievedExpression.Status)
	}
}

func TestSetResultExpression(t *testing.T) {

	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	// Создаем пользователя для получения creatorID
	userForm := UserCreateForm{
		Username: "testuser",
		Password: "password123",
	}
	user, err := db.CreateUser(userForm)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	form := CalculateRequest{
		Expression: "2+2",
	}
	expression, err := db.CreateExpression(user.ID, form)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = db.SetResultExpression(expression.ID, "completed", 4.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updatedExpression, err := db.GetExpressionByID(expression.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updatedExpression.Status != "completed" {
		t.Errorf("expected status 'completed', got %s", updatedExpression.Status)
	}
	if updatedExpression.Result != 4.0 {
		t.Errorf("expected result 4.0, got %f", updatedExpression.Result)
	}
}

func TestGetAllExpressions(t *testing.T) {

	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	// Создаем пользователя для получения creatorID
	userForm := UserCreateForm{
		Username: "testuser",
		Password: "password123",
	}
	user, err := db.CreateUser(userForm)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	expressions := []CalculateRequest{
		{Expression: "2+2"},
		{Expression: "3+3"},
	}
	for _, form := range expressions {
		_, err := db.CreateExpression(user.ID, form)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	allExpressions, err := db.GetAllExpressions()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(allExpressions) != len(expressions) {
		t.Errorf("expected %d expressions, got %d", len(expressions), len(allExpressions))
	}
}

func TestGetAllExpressionsByUserID(t *testing.T) {

	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	userID := "test-user-id"
	expressions := []struct {
		ID         string
		Expression string
		Status     string
		Result     float64
	}{
		{ID: "1", Expression: "2+2", Status: "completed", Result: 4.0},
		{ID: "2", Expression: "3+3", Status: "completed", Result: 6.0},
	}
	for _, expr := range expressions {
		_, err := db.dbConnection.Exec(
			"INSERT INTO expressions (id, expression, status, result, creator_id) VALUES (?, ?, ?, ?, ?)",
			expr.ID, expr.Expression, expr.Status, expr.Result, userID,
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	userExpressions, err := db.GetAllExpressionsByUserID(userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(userExpressions) != len(expressions) {
		t.Errorf("expected %d expressions, got %d", len(expressions), len(userExpressions))
	}
}

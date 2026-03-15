package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/user/todo-api/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	// Use in-memory database for testing
	var err error
	db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Todo{})

	os.Exit(m.Run())
}

func TestSwaggerVersionAndSchemaSync(t *testing.T) {
	// Read the generated swagger.json
	data, err := os.ReadFile("docs/swagger.json")
	require.NoError(t, err, "swagger.json should exist. Run 'make swag-init' if it's missing.")

	var swagger struct {
		Info struct {
			Version string `json:"version"`
		} `json:"info"`
		Definitions json.RawMessage `json:"definitions"`
	}

	err = json.Unmarshal(data, &swagger)
	require.NoError(t, err)

	// Current expected version after adding due_date
	expectedVersion := "1.1.0"

	// 1. Check if the version is updated in main.go/docs
	assert.Equal(t, expectedVersion, swagger.Info.Version, "The API version should be %s", expectedVersion)

	// 2. Hash the definitions to detect any schema changes
	hash := sha256.New()
	hash.Write(swagger.Definitions)
	definitionsHash := hex.EncodeToString(hash.Sum(nil))

	// This is the hash of the 'definitions' at version 1.1.0
	// If you change the models, this hash will change, forcing you to
	// acknowledge the change by updating it here AND incrementing the version.
	expectedHash := "6fb6e2c8118955c8a12f17e02291df6051cdc1e1c28f2a192da0225d7f7cea91"

	if definitionsHash != expectedHash {
		t.Errorf("Schema (definitions) has changed but expected hash remains %s.\nNew hash: %s\n\nIf you've updated the schema, you MUST:\n1. Increment @version in main.go\n2. Run 'make swag-init'\n3. Update expectedVersion and expectedHash in main_test.go", expectedHash, definitionsHash)
	}
}

func TestGetTodos(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/todos", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateTodo(t *testing.T) {
	router := setupRouter()

	todo := models.CreateTodoInput{
		Title:   "Test Todo",
		DueDate: time.Now(),
	}
	jsonValue, _ := json.Marshal(todo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/todos", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Todo
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Test Todo", response.Title)
}

func TestPing(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "pong", response["message"])
}

func TestPingV1(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "pong from v1", response["message"])
}

//go:generate swag init
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/user/todo-api/docs"
	"github.com/user/todo-api/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// @title           Todo List API
// @version         1.0
// @description     A simple todo list API in Golang using Gin and GORM.
// @host            localhost:8080
// @BasePath        /api/v1

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("todos.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.Todo{})
	if err != nil {
		log.Fatal("failed to migrate database")
	}

	r := setupRouter()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so no need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173", "https://jerri-azeotropic-bertha.ngrok-free.dev"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	v1 := r.Group("/api/v1")
	{
		todos := v1.Group("/todos")
		{
			todos.GET("", getTodos)
			todos.POST("", createTodo)
			todos.GET("/:id", getTodo)
			todos.PUT("/:id", updateTodo)
			todos.DELETE("/:id", deleteTodo)
		}
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

// getTodos godoc
// @Summary      List todos
// @Description  get todos
// @Tags         todos
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Todo
// @Router       /todos [get]
func getTodos(c *gin.Context) {
	var todos []models.Todo
	db.Find(&todos)
	c.JSON(http.StatusOK, todos)
}

// createTodo godoc
// @Summary      Create a todo
// @Description  create a new todo item
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        todo  body      models.Todo  true  "Todo object"
// @Success      201  {object}  models.Todo
// @Router       /todos [post]
func createTodo(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&todo)
	c.JSON(http.StatusCreated, todo)
}

// getTodo godoc
// @Summary      Get a todo
// @Description  get a todo item by ID
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Todo ID"
// @Success      200  {object}  models.Todo
// @Router       /todos/{id} [get]
func getTodo(c *gin.Context) {
	var todo models.Todo
	if err := db.First(&todo, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	c.JSON(http.StatusOK, todo)
}

// updateTodo godoc
// @Summary      Update a todo
// @Description  update an existing todo item
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        id    path      int          true  "Todo ID"
// @Param        todo  body      models.Todo  true  "Updated todo object"
// @Success      200   {object}  models.Todo
// @Router       /todos/{id} [put]
func updateTodo(c *gin.Context) {
	var todo models.Todo
	if err := db.First(&todo, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	var input models.Todo
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&todo).Updates(input)
	c.JSON(http.StatusOK, todo)
}

// deleteTodo godoc
// @Summary      Delete a todo
// @Description  delete a todo item by ID
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Todo ID"
// @Success      204  "No Content"
// @Router       /todos/{id} [delete]
func deleteTodo(c *gin.Context) {
	var todo models.Todo
	if err := db.First(&todo, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	db.Delete(&todo)
	c.Status(http.StatusNoContent)
}

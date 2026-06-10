package main

import (
	"log"
	"net/http"

	"task/internal/auth"
	"task/internal/database"
	"task/internal/middleware"
	"task/internal/tasks"

	"github.com/joho/godotenv"
	_ "task/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Task Management API
// @version 1.0
// @description This is a sample task management server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	database.InitDB()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth/register", auth.RegisterHandler)
	mux.HandleFunc("POST /auth/login", auth.LoginHandler)

	taskHandlers := http.NewServeMux()
	taskHandlers.HandleFunc("POST /tasks", tasks.CreateTaskHandler)
	taskHandlers.HandleFunc("GET /tasks", tasks.GetAllTasksHandler)
	taskHandlers.HandleFunc("GET /tasks/{id}", tasks.GetTaskHandler)
	taskHandlers.HandleFunc("PUT /tasks/{id}", tasks.UpdateTaskHandler)
	taskHandlers.HandleFunc("DELETE /tasks/{id}", tasks.DeleteTaskHandler)

	mux.Handle("/tasks", middleware.AuthMiddleware(taskHandlers))
	mux.Handle("/tasks/", middleware.AuthMiddleware(taskHandlers))

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

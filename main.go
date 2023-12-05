package main

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var tasks []Task

func main() {
	r := gin.Default()

	r.GET("/tasks", BasicAuth, getTasks)
	r.POST("/tasks", addTask)
	r.PUT("/tasks/:id", updateTask)
	r.DELETE("/tasks/:id", deleteTask)

	r.Run(":8080")
}

func getTasks(c *gin.Context) {
	if len(tasks) == 0 {
		response := gin.H{
			"href":    "/tasks",
			"message": "data not found",
			"status":  200,
		}
		c.JSON(http.StatusOK, response)
		return
	}

	response := struct {
		Message string `json:"message"`
		Data    []Task `json:"data"`
		Status  string `json:"status"`
	}{
		Message: "Success",
		Data:    tasks,
		Status:  "200",
	}

	c.JSON(http.StatusOK, response)
}

func addTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBind(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tasks = append(tasks, task)
	c.JSON(http.StatusCreated, task)
}

func updateTask(c *gin.Context) {
	taskID := c.Param("id")

	// Cari task berdasarkan ID
	var foundTask *Task
	for i := range tasks {
		if tasks[i].ID == taskID {
			foundTask = &tasks[i]
			break
		}
	}

	// Jika task tidak ditemukan
	if foundTask == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Bind data yang baru dari request
	var updatedTask Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Perbarui task
	foundTask.Title = updatedTask.Title
	foundTask.Content = updatedTask.Content

	c.JSON(http.StatusOK, foundTask)
}

func deleteTask(c *gin.Context) {
	taskID := c.Param("id")
	for i, task := range tasks {
		if task.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			c.JSON(http.StatusCreated, nil)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

func BasicAuth(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")

	if !strings.Contains(authHeader, "Basic") {
		result := gin.H{
			"status":  http.StatusForbidden,
			"message": "invalid token",
			"href":    c.Request.RequestURI,
		}
		c.JSON(http.StatusForbidden, result)
		c.Abort()
		return
	}

	clientSecret := "1jutadolar2024"
	clientID := "bikin.dev"
	tokenString := strings.Replace(authHeader, "Basic ", "", -1)
	myToken := clientID + ":" + clientSecret
	myBasicAuth := base64.StdEncoding.EncodeToString([]byte(myToken))
	if tokenString != myBasicAuth {
		result := gin.H{
			"status":  http.StatusUnauthorized,
			"message": "invalid authentication",
			"href":    c.Request.RequestURI,
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
		return
	}
}

package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// task represents data about a task
type task struct {
	TaskID      string `json:"task_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
}

// tasks slice to seed task data
var tasks = []task{
	{TaskID: "1", Title: "Buy groceries", Description: "Go to the market and buy all weeks groceries", IsCompleted: false},
	{TaskID: "2", Title: "Pancakes", Description: "Prepare pancakes using EX home baking flour", IsCompleted: false},
	{TaskID: "3", Title: "Chapter 7 Interfaces", Description: "Read the whole of chapter 7 of idiomatic Go interfaces", IsCompleted: false},
}

// mu guards tasks. It must be shared by every handler.
var mu sync.RWMutex

func main() {
	router := gin.Default()
	router.GET("/tasks/:id", getTaskByID)
	router.GET("/tasks", getAllTasks)
	router.POST("/tasks", postTask)

	router.Run("localhost:8080")
}

// getAllTasks responds with the list of all tasks as JSON
func getAllTasks(c *gin.Context) {
	mu.RLock()
	snapshot := make([]task, len(tasks))
	copy(snapshot, tasks)
	mu.RUnlock()

	c.IndentedJSON(http.StatusOK, snapshot)
}

// getTaskByID locates the task whose ID matches the id parameter
func getTaskByID(c *gin.Context) {
	id := c.Param("id")

	mu.RLock()
	defer mu.RUnlock()

	for _, t := range tasks {
		if t.TaskID == id {
			c.IndentedJSON(http.StatusOK, t)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
}

// postTask adds a task from the JSON in the request body
func postTask(c *gin.Context) {
	var newTask task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	tasks = append(tasks, newTask)
	mu.Unlock()

	c.IndentedJSON(http.StatusCreated, newTask)
}
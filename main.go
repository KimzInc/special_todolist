package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//server holds the dependencies the handlers need. It depends on the
// TaskStore interface, never on a concrete type.
type server struct{
	store TaskStore
}

func main() {
	srv := &server{store: NewMemoryStore()}

	router := gin.Default()
	router.GET("/tasks/:id", srv.getTaskByID)
	router.GET("/tasks", srv.getAllTasks)
	router.POST("/tasks", srv.postTask)

	if err := router.Run("localhost:8080"); err != nil {
		log.Fatal(err)
	}
}


// getAllTask responds with the list of all tasks as JSON.
func (s *server) getAllTasks(c *gin.Context){
	tasks, err := s.store.List(c.Request.Context())

	if err != nil {
		log.Printf("list tasks: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.IndentedJSON(http.StatusOK, tasks)
}


// getTaskByID responds with the task whose ID matches the id parameter.
func (s *server) getTaskByID(c *gin.Context) {
	t, err := s.store.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
			return
		}
		log.Printf("get task: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.IndentedJSON(http.StatusOK, t)
}
 
// postTask adds a task from the JSON in the request body.
func (s *server) postTask(c *gin.Context) {
	var newTask Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
 
	if err := s.store.Add(c.Request.Context(), newTask); err != nil {
		if errors.Is(err, ErrDuplicate) {
			c.IndentedJSON(http.StatusConflict, gin.H{"error": "task_id already exists"})
			return
		}
		log.Printf("add task: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.IndentedJSON(http.StatusCreated, newTask)
}
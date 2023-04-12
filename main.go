package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type task struct {
	Task        string `json:"task"`
	Tag         string `json:"tag"`
	Goal        string `json:"goal"`
	Color       string `json:"color"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}

type taskResponse struct {
	Task        string
	Tag         string
	Goal        string
	Color       string
	Icon        string
	Description string
}

var tasks = []task{
	{Task: "Reading", Tag: "Productivity", Goal: "Reading", Color: "yellow", Icon: "book", Description: "30 mins of mindful reading"},
	{Task: "Coding", Tag: "Productivity", Goal: "Programming", Color: "blue", Icon: "computer", Description: "30 mins of programming"},
}

func getTasks(c *gin.Context) {

	// TODO: Why is SSL Mode disabled
	connStr := "postgres://postgres:password@localhost/atracker?sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		fmt.Println(err.Error())
	}
	tasks := []taskResponse{}

	rows, err := db.Query(`select * from task`)

	if err != nil {
		c.IndentedJSON(http.StatusBadGateway, err.Error())
	}

	for rows.Next() {
		var task taskResponse

		err := rows.Scan(&task.Task, &task.Tag, &task.Goal,
			&task.Color, &task.Icon, &task.Description)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}

		tasks = append(tasks, task)
	}

	rows.Close()
	db.Close()

	c.IndentedJSON(http.StatusOK, tasks)

}

func createTask(c *gin.Context) {

	connStr := "postgres://postgres:password@localhost/atracker?sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		fmt.Println("Error")
	}

	var newTask task

	if err := c.BindJSON(&newTask); err != nil {
		return
	}

	fmt.Println(newTask)

	result, err := db.Exec(`insert into task (
		task, tag, goal,
		color, icon, description
	) values (
		$1, %2, $3,
		$4, %5, $6);`, newTask.Task, newTask.Tag, newTask.Goal,
		newTask.Icon, newTask.Color, newTask.Description)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	}

	count, err := result.RowsAffected()

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	}

	db.Close()

	c.IndentedJSON(http.StatusCreated, gin.H{"Rows affected": count})

}

func getTask(c *gin.Context) {
	connStr := "postgres://postgres:password@localhost/atracker?sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		fmt.Println("Error")
	} else {
		fmt.Println("Database connected")
	}

	db.Close()

	fmt.Println("Database connection closed")

	task := c.Param("id")

	for _, a := range tasks {
		if a.Task == task {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Task not found"})
}

func main() {

	router := gin.Default()
	router.GET("/tasks", getTasks)
	router.POST("/task", createTask)
	router.GET("/task/:id", getTask)

	router.Run("localhost:8080")
}

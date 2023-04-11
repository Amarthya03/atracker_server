package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type task struct {
	Task        string `json:"task"`
	Tag         string `json:"tag"`
	Goal        string `json:"goal"`
	Color       string `json:"color"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}

var tasks = []task{
	{Task: "Reading", Tag: "Productivity", Goal: "Reading", Color: "yellow", Icon: "book", Description: "30 mins of mindful reading"},
	{Task: "Coding", Tag: "Productivity", Goal: "Programming", Color: "blue", Icon: "computer", Description: "30 mins of programming"},
}

func getTasks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, tasks)
}

func main() {
	router := gin.Default()
	router.GET("/tasks", getTasks)

	router.Run("localhost:8080")
}

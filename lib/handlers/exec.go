package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liron-navon/code-runner/lib/internals"
)

type CodeExecutionRequest struct {
	Code     string `json:"code" binding:"required"`
	Language string `json:"language" binding:"required"`
}

func HandleExec(c *gin.Context) {
	var request CodeExecutionRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// get the language context
	langContext := internals.ProgrammingLanguages[request.Language]
	if len(langContext.Name) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("Unknown language %s", request.Language),
		})
		return
	}

	output, err := langContext.SafeExecuteCode(request.Code)
	if err != nil {
		log.Println("failed running go:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"output": output,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"output": output,
	})
}

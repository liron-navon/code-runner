package internals

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type job struct {
	ctx *CodeExecutionContext
	code string
}

var jobs = make(chan job, 500)
var isReadyChannel = make(chan bool, 1)
var didInitiateIsReadyChannel = false;

// log steps
const verbose = false

// the name where we store the temporary files (auto deleted when done with)
const temporaryDirectoryName = "tmp"

// the file name passed to docker, for example in java it will be "main.java" and in go "main.go"
const dockerCodeFileName = "main"

// SafeExecuteCode -  executes the code from a buffered queue, to make sure only one OP is running at a time
func (ctx *CodeExecutionContext) SafeExecuteCode(code string) ([]string, error) {
	// make sure we initiate the queue once
	if !didInitiateIsReadyChannel {
		didInitiateIsReadyChannel = true;
		isReadyChannel <- true
	}

	// create a new code execution job
	jobs <- job{
		ctx:  ctx,
		code: code,
	}

	// wait for the queue to be ready
	_ = <- isReadyChannel

	// wait for the next job to come in
	nextJob := <- jobs

	// execute the job
	lines, err :=  nextJob.ctx.ExecuteCode(nextJob.code)

	// notify that the queue is ready
	isReadyChannel <- true

	return lines, err
}

// ExecuteCode - accepts a context and a piece of code to execute,
// and returns the standard output lines for the execiution
func (ctx *CodeExecutionContext) ExecuteCode(code string) ([]string, error) {
	if verbose {
		fmt.Println("--- Input ---")
		fmt.Println(code)
	}

	// make sure the file is always an absolute path
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// create the file path

	uid := uuid.New()
	pathToDirectory := filepath.Join(currentDir, temporaryDirectoryName, ctx.Name, uid.String())
	fileName := fmt.Sprintf("%s.%s", dockerCodeFileName, ctx.CodeFilePostfix)
	pathToFile := filepath.Join(pathToDirectory, fileName)
	pathToDirectoryInDocker := fmt.Sprintf("/%s/%s", temporaryDirectoryName, ctx.Name)
	pathToFileInContainer := fmt.Sprintf("%s/%s", pathToDirectoryInDocker, fileName)

	// make sure the directory exists
	err = ValidateDirectory(pathToDirectory)
	if err != nil {
		return nil, err
	}

	// write the file
	err = WriteFileSync(pathToFile, code)
	if err != nil {
		return nil, err
	}

	// prepare and execute the command
	codeRunnerCommand := fmt.Sprintf(ctx.CodeRunnerCommand, pathToFileInContainer)
	commandToExecute := fmt.Sprintf("docker run -v %s:%s -t %s %s", pathToDirectory, pathToDirectoryInDocker, ctx.Image, codeRunnerCommand)

	if verbose {
		fmt.Println("--- Docker Execution command ---")
		fmt.Println(commandToExecute)
	}

	// execute command and split into lines
	out, err := ExecuteCommandline(time.Second*time.Duration(ctx.MaxSeconds), commandToExecute, []string{})
	lines := cleanOutputLines(strings.Split(string(out), "\n"))
	if err != nil {
		return lines, err
	}

	err = DeleteFile(pathToDirectory)
	if err != nil {
		return nil, err
	}

	if verbose {
		fmt.Println("--- Output ---")
		for i, line := range lines {
			fmt.Println(i, ":", line)
		}
	}

	return lines, nil
}

// cleanOutputLines - trims lines and remove empty EOF lines
func cleanOutputLines(lines []string) []string {
	var newLines []string
	for _, item := range lines {
		if item != "" {
			newLines = append(newLines, strings.TrimSpace(item))
		}
	}
	return newLines
}

//CodeExecutionContext - struct
type CodeExecutionContext struct {
	Name              string
	MaxSeconds        int
	CodeRunnerCommand string
	CodeFilePostfix   string
	Image             string
	NextStage         *CodeExecutionContext
}

// ProgrammingLanguages - contexts for supported programming languages
var ProgrammingLanguages = map[string]CodeExecutionContext{
	"go": CodeExecutionContext{
		Name:              "go",
		CodeFilePostfix:   "go",
		CodeRunnerCommand: "go run %s",
		Image:             "golang:1.12",
		MaxSeconds:        10,
	},
	"node": CodeExecutionContext{
		Name:              "node",
		CodeFilePostfix:   "javascript",
		CodeRunnerCommand: "node %s",
		Image:             "node:10",
		MaxSeconds:        10,
	},
	"python3": CodeExecutionContext{
		Name:              "python3",
		CodeFilePostfix:   "py",
		CodeRunnerCommand: "python %s",
		Image:             "python:3.7",
		MaxSeconds:        10,
	},
	"java": CodeExecutionContext{
		Name:              "java",
		CodeFilePostfix:   "java",
		CodeRunnerCommand: "java %s",
		Image:             "openjdk:11",
		MaxSeconds:        10,
	},
}

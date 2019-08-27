package internals

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const isLinux = runtime.GOOS == "linux"

func DeleteFile(path string) error {
	return os.RemoveAll(path)
}

func ValidateDirectory(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func WriteFileSync(path string, data string) error {
	_, err := os.Create(path)
	if err != nil {
		return err
	}

	// open file using READ & WRITE permission
	file, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	// write some text line-by-line to file
	_, err = file.WriteString(data)
	if err != nil {
		return err
	}
	// save changes
	return file.Sync()
}

// ExecuteCommandline - executes a relatively "safe" command line, with a timeout
func ExecuteCommandline(time time.Duration, command string, extraArgs []string) (string, error) {
	// Create a new context and add a timeout to it
	ctx, cancel := context.WithTimeout(context.Background(), time)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	// Create the command with our context
	args := strings.Split(command, " ")
	var cmd *exec.Cmd

	if len(args) == 1 {
		cmd = exec.CommandContext(ctx, args[0], extraArgs...)
	} else {
		cmd = exec.CommandContext(ctx, args[0], append(args[1:], extraArgs...)...)
	}

	cmd.Wait()
	// This time we can simply use Output() to get the result.
	out, err := cmd.CombinedOutput()

	// We want to check the context error to see if the timeout was executed.
	// The error returned by cmd.Output() will be OS specific based on what
	// happens when a process is killed.
	if ctx.Err() == context.DeadlineExceeded {
		return "", ctx.Err()
	}

	// If there's no context error, we know the command completed (or errored).
	if err != nil {
		return string(out), err
	}
	return string(out), nil
}

// findFileWithExtention - finds a file with a specified extention
func findFileWithExtention(dir string, ext string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == fmt.Sprintf(".%s", ext) {
				return file.Name(), nil
			}
		}
	}
	return "", fmt.Errorf("No file with extention %s found", ext)
}

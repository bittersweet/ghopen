package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func getCurrentBranch() string {
	branch, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if nerr, ok := err.(*exec.ExitError); ok {
		log.Fatalf("Exit error: %s", nerr.Error())
	}

	// Trim \n character from rev-parse output
	branch = branch[:len(branch)-1]
	return string(branch)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		return false
	}

	return true
}

func main() {
	url, err := exec.Command("git", "config", "remote.origin.url").Output()
	if nerr, ok := err.(*exec.ExitError); ok {
		log.Fatalf("Exit error: %s", nerr.Error())
	}

	r := regexp.MustCompile(`github.com:(.+)\.git`)

	repository := r.FindAllStringSubmatch(string(url), -1)
	if repository == nil {
		log.Fatalf("No repository found")
	}
	fullUrl := fmt.Sprintf("https://github.com/%s", repository[0][1])

	gitRoot, err := exec.Command("git", "rev-parse", "--sq", "--show-toplevel").Output()
	if nerr, ok := err.(*exec.ExitError); ok {
		log.Fatalf("Exit error: %s", nerr.Error())
	}
	// Trim \n character from rev-parse output
	gitRoot = gitRoot[:len(gitRoot)-1]
	gitDir := string(gitRoot)

	pwd := os.Getenv("PWD")
	branch := getCurrentBranch()

	if len(os.Args) > 1 {
		filename := os.Args[1]
		if !fileExists(filename) {
			fmt.Printf("File does not exist: %s\n", filename)
			os.Exit(1)
		}

		splitted := strings.SplitAfter(pwd, gitDir)

		fullUrl = fmt.Sprintf("%s/tree/%s%s", fullUrl, branch, splitted[1])
		fullUrl = fmt.Sprintf("%s/%s", fullUrl, filename)
	} else if pwd != gitDir {
		splitted := strings.SplitAfter(pwd, gitDir)

		fullUrl = fmt.Sprintf("%s/tree/%s%s", fullUrl, branch, splitted[1])
	}

	fmt.Printf("Opening url: %s\n", fullUrl)
	exec.Command("open", fullUrl).Run()
}

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

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

	if len(os.Args) > 1 {
		splitted := strings.SplitAfter(pwd, gitDir)

		branch := "master" // TODO: update to get current git branch
		fullUrl = fmt.Sprintf("%s/tree/%s%s", fullUrl, branch, splitted[1])
		fullUrl = fmt.Sprintf("%s/%s", fullUrl, os.Args[1])
	} else if pwd != gitDir {
		// index 1 will contain the difference between pwd and git root
		splitted := strings.SplitAfter(pwd, gitDir)

		branch := "master" // TODO: update to get current git branch
		fullUrl = fmt.Sprintf("%s/tree/%s%s", fullUrl, branch, splitted[1])
	}

	fmt.Printf("Opening url: %s\n", fullUrl)
	exec.Command("open", fullUrl).Run()
}

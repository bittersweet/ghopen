package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func gitCommand(args ...string) string {
	output, err := exec.Command("git", args...).Output()
	if _, ok := err.(*exec.ExitError); ok {
		fmt.Println("GHOpen error: Received exit status 1 from git")
		os.Exit(1)
	}

	// Trim \n character from rev-parse output
	return string(output[:len(output)-1])
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		return false
	}

	return true
}

func getRepository() string {
	url := gitCommand("config", "remote.origin.url")
	r := regexp.MustCompile(`github.com:(.+)\.git`)
	repository := r.FindAllStringSubmatch(string(url), -1)
	if repository == nil {
		fmt.Println("GHOpen error: no remote repository found")
		os.Exit(1)
	}
	return fmt.Sprintf("https://github.com/%s", repository[0][1])
}

func getGitRoot() string {
	return gitCommand("rev-parse", "--sq", "--show-toplevel")
}

func main() {
	fullUrl := getRepository()
	gitRoot := getGitRoot()

	pwd := os.Getenv("PWD")
	sha := gitCommand("rev-parse", "HEAD")

	if len(os.Args) > 1 {
		filename := os.Args[1]
		if !fileExists(filename) {
			fmt.Printf("File does not exist: %s\n", filename)
			os.Exit(1)
		}

		splitted := strings.SplitAfter(pwd, gitRoot)

		fullUrl = fmt.Sprintf("%s/tree/%s%s", fullUrl, sha, splitted[1])
		fullUrl = fmt.Sprintf("%s/%s", fullUrl, filename)

		// argument 2 contains the line number
		if len(os.Args) > 2 {
			fullUrl = fmt.Sprintf("%s#L%s", fullUrl, os.Args[2])
		}
	} else if pwd != gitRoot {
		splitted := strings.SplitAfter(pwd, gitRoot)

		fullUrl = fmt.Sprintf("%s/tree/%s%s", fullUrl, sha, splitted[1])
	}

	fmt.Printf("Opening url: %s\n", fullUrl)
	exec.Command("open", fullUrl).Run()
}

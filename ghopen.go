package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
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

	commitSha := flag.String("commit", "default", "sha1")
	flag.Parse()

	if *commitSha != "default" {
		fullUrl = fmt.Sprintf("%s/commit/%s", fullUrl, *commitSha)
	} else if len(os.Args) > 1 {
		filename := os.Args[1]
		if !fileExists(filename) {
			fmt.Printf("File does not exist: %s\n", filename)
			os.Exit(1)
		}

		revOutput := gitCommand("rev-parse", "--show-toplevel", filename)
		revOutputSplit := strings.Split(revOutput, "\n")

		// /home/user/project
		gitRoot := revOutputSplit[0]
		// /home/user/project/directory/file.extension
		fullFilePath := revOutputSplit[1]

		relativePath := strings.TrimPrefix(fullFilePath, gitRoot)
		relativePath = strings.TrimPrefix(relativePath, "/")

		fullUrl = fmt.Sprintf("%s/blob/%s/%s", fullUrl, sha, relativePath)

		// argument 2 contains the line number
		if len(os.Args) > 2 {
			fullUrl = fmt.Sprintf("%s#L%s", fullUrl, os.Args[2])
		}
	} else if pwd != gitRoot {
		splitted := strings.SplitAfter(pwd, gitRoot)

		fullUrl = fmt.Sprintf("%s/tree/%s%s", fullUrl, sha, splitted[1])
	}

	fmt.Printf("Opening url: %s\n", fullUrl)
	if runtime.GOOS == "linux" {
		exec.Command("xdg-open", fullUrl).Run()
	} else {
		exec.Command("open", fullUrl).Run()
	}
}

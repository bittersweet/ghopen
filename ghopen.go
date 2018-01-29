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

	gitDir, err := exec.Command("git", "rev-parse", "--sq", "--show-toplevel").Output()
	if nerr, ok := err.(*exec.ExitError); ok {
		log.Fatalf("Exit error: %s", nerr.Error())
	}
	dirString := string(gitDir)
	dirString = strings.TrimRight(dirString, "\n")

	// check pwd with gitdir
	pwd := os.Getenv("PWD")
	if pwd == dirString {
		// exactly the same, we can open the root page on GH
		fmt.Println("pwd == gitdir")
	} else {
		// index 1 will contain the difference between pwd and git root
		splitted := strings.SplitAfter(pwd, dirString)
		fmt.Printf("%#v\n", splitted)

		branch := "master" // TODO: update to get current git branch
		fullUrl = fmt.Sprintf("%s/tree/%s%s", fullUrl, branch, splitted[1])
	}

	fmt.Printf("Opening url: %s\n", fullUrl)
	exec.Command("open", fullUrl).Run()
}

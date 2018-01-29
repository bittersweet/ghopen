package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
)

func main() {
	// walk up to we find a .git ?
	// or just run git commands in current directory
	// cmd := "/usr/local/bin/git config --get remote.origin.url"
	url, err := exec.Command("git", "config", "remote.origin.url").Output()
	if nerr, ok := err.(*exec.ExitError); ok {
		log.Fatalf("Exit error: %s", nerr.Error())
	}

	r := regexp.MustCompile(`github.com:(.+)\.git`)

	repository := r.FindAllStringSubmatch(string(url), -1)
	if repository == nil {
		log.Fatalf("No repository found")
	}
	// check if current directory contains the .git dir
	fullUrl := fmt.Sprintf("https://github.com/%s", repository[0][1])

	fmt.Printf("Opening url: %s\n", fullUrl)
	exec.Command("open", fullUrl).Run()
}

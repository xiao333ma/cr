package main

import (
	"os/exec"
	"strings"
)

var crPath = "/merge_requests/new"

func main() {

	c := exec.Command("git", "ls-remote", "--get-url")
	remote, _ := c.CombinedOutput()

	/**
	url 有两种

	git@github.com:xiao333ma/Friday_Server.git
	https://github.com/xiao333ma/Friday_Server.git
	*/

	stringArr := strings.Split(string(remote), "@")

	host := stringArr[len(stringArr)-1]

	stringArr = strings.Split(host, ".git")

	host = stringArr[0]

	if !strings.HasPrefix(host, "https://") {
		host = strings.ReplaceAll(host, ":", "/")
		host = "https://" + host
	}

	host += crPath
	c = exec.Command("open", host)
	c.Run()
}

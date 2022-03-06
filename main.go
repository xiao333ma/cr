package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var crPath = "/merge_requests/new"

var f = flag.Bool("f", false, "发起一个 current branch ➜ feature 的 CR")
var d = flag.Bool("d", false, "发起一个 current branch ➜ develop 的 CR")
var s = flag.String("s", "", "source branch")
var t = flag.String("t", "", "target branch")
var p = flag.String("p", "", "子目录")

var reset = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"

func main() {

	flag.Parse()

	if len(os.Args) == 1 {
		open()
		return
	}

	if len(*p) > 0 && !enterTargetPath(*p) {
		return
	}

	if *f {
		mergeToFeature()
		return
	}

	if *d {
		mergeToDevelop()
		return
	}

	if len(*s) > 0 && len(*t) > 0 {
		merge(*s, *t)
		return
	}
}

func enterTargetPath(path string) bool {
	pwd, _ := os.Getwd()
	err := os.Chdir(pwd + "/" + *p)
	if err != nil {
		fmt.Println("无法进入", *p)
		return false
	}
	return true
}

func open() {
	openURL(getNewCRPath())
}

func mergeToFeature() {

	currentBranch := getCurrentBranch()
	featureBranch := getFeatureBranch(currentBranch)

	if !isRemoteBranchExist(featureBranch, getRepoURL()) {
		fmt.Println(red, "无法发起 CR:", currentBranch, "对应的 feature 分支", featureBranch, "不存在", reset)
		return
	}
	url := buildMergeRequestURL(currentBranch, featureBranch)
	openURL(url)
}

func mergeToDevelop() {

	currentBranch := getCurrentBranch()
	url := buildMergeRequestURL(currentBranch, "develop")
	openURL(url)
}

func merge(sourceBranch string, targetBranch string) {

	if !isRemoteBranchExist(sourceBranch, getRepoURL()) {
		fmt.Println(red, "无法发起 CR:", sourceBranch, "不存在", reset)
		return
	}

	if !isRemoteBranchExist(targetBranch, getRepoURL()) {
		fmt.Println(red, "无法发起 CR:", targetBranch, "不存在", reset)
		return
	}

	url := buildMergeRequestURL(sourceBranch, targetBranch)
	openURL(url)
}

func buildMergeRequestURL(sourceBranch string, targetBranch string) string {
	url := getNewCRPath()

	url += "?merge_request[source_branch]=" + sourceBranch
	url += "&"
	url += "merge_request[target_branch]=" + targetBranch
	fmt.Println(green, "🍺", sourceBranch, "➜", targetBranch, reset)

	return url
}

func openURL(url string) {
	c := exec.Command("open", url)
	c.Run()
}

func getNewCRPath() string {
	url := getRepoURL()
	stringArr := strings.Split(url, ".git")
	url = stringArr[0]
	url += crPath
	return url
}

func getRepoURL() string {

	c := exec.Command("git", "ls-remote", "--get-url")
	remote, _ := c.CombinedOutput()
	stringArr := strings.Split(string(remote), "@")
	host := stringArr[len(stringArr)-1]
	if !strings.HasPrefix(host, "https://") && !strings.HasPrefix(host, "http://") {
		host = strings.ReplaceAll(host, ":", "/")
		host = "https://" + host
	}
	return strings.TrimSuffix(host, "\n")
}

func getCurrentBranch() string {
	c := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	remote, _ := c.CombinedOutput()
	return strings.TrimSuffix(string(remote), "\n")
}

func getFeatureBranch(branchName string) string {

	arr := strings.Split(branchName, "/")
	name := arr[len(arr)-1]
	return "feature/" + name
}

func isRemoteBranchExist(branchName string, repoURL string) bool {
	c := exec.Command("git", "ls-remote", "--heads", repoURL, branchName)
	remote, _ := c.CombinedOutput()
	return len(remote) > 0
}

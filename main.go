package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var crPath = "/merge_requests"

var f = flag.Bool("f", false, "发起一个 current branch ➜ feature 的 CR")
var d = flag.Bool("d", false, "发起一个 current branch ➜ develop 的 CR")
var r = flag.Bool("r", false, "发起一个 current branch ➜ release 的 CR")
var m = flag.Bool("m", false, "发起一个 current branch ➜ master 的 CR")
var M = flag.Bool("M", false, "发起一个 current branch ➜ Main 的 CR")
var l = flag.Bool("l", false, "打开当前 repo 的 CR list")
var s = flag.String("s", "", "source branch, 配合 -t 使用 可发起 source ➜ target 的 CR")
var t = flag.String("t", "", "target branch, 配合 -s 使用 可发起 source ➜ target 的 CR")
var p = flag.String("p", "", "子目录，进入子目录发起 CR，省去了 cd 命令")

var reset = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"

func main() {

	flag.Parse()

	if len(*p) > 0 && !enterTargetPath(*p) {
		return
	}

	if !isGitRepo() {
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

	if *r {
		mergeToRelease()
		return
	}

	if *m {
		mergeToMaster()
		return
	}

	if *M {
		mergeToMain()
		return
	}

	if len(*s) > 0 && len(*t) > 0 {
		merge(*s, *t)
		return
	}

	if *l {
		openURL(getCRBaseURL())
		return
	}

	open()
}

func enterTargetPath(path string) bool {
	targetPath := ""
	if strings.HasPrefix(path, "/") {
		targetPath = path
	} else {
		pwd, _ := os.Getwd()
		targetPath = pwd + "/" + path
	}

	err := os.Chdir(targetPath)
	if err != nil {
		fmt.Println(red, err, reset)
		return false
	}
	return true
}

func open() {
	openURL(getCRNewURL())
}

func mergeToFeature() {

	currentBranch := getCurrentBranch()
	featureBranch := getFeatureBranch(currentBranch)

	merge(currentBranch, featureBranch)
}

func mergeToDevelop() {
	merge(getCurrentBranch(), "develop")
}

func mergeToMaster()  {
	merge(getCurrentBranch(), "master")

}

func mergeToMain()  {
	merge(getCurrentBranch(), "main")
}

func mergeToRelease()  {

	branches , err := findReleaseBranches()
	if err != nil {
		fmt.Println(red, "无法发起 CR: 找不到对应的 release 分支", reset)
		return
	}

	releaseBranchName := getTargetReleaseBranch(branches)
	if len(releaseBranchName) > 0 {
		merge(getCurrentBranch(), releaseBranchName)
	}
}

func merge(sourceBranch string, targetBranch string) {

	if !strings.Contains(sourceBranch, "release") {
		if !isRemoteBranchExist(sourceBranch, getRepoGitURL()) {
			fmt.Println(red, "无法发起 CR:", sourceBranch, "不存在", reset)
			return
		}
	}

	if !strings.Contains(targetBranch, "release") {
		if !isRemoteBranchExist(targetBranch, getRepoGitURL()) {
			fmt.Println(red, "无法发起 CR:", targetBranch, "不存在", reset)
			return
		}
	}

	url := buildMergeRequestURL(sourceBranch, targetBranch)
	openURL(url)
}

func buildMergeRequestURL(sourceBranch string, targetBranch string) string {
	url := getCRNewURL()

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

func getCRBaseURL() string {
	return getRepoURL() + crPath
}

func getCRNewURL() string {
	url := getCRBaseURL()
	url += "/new"
	return url
}

/*返回当前 repo 的 URL */
func getRepoURL() string {
	url := getRepoGitURL()
	stringArr := strings.Split(url, ".git")
	url = stringArr[0]
	return url
}

/*返回当前 repo 的 URL，后边带 .git */
func getRepoGitURL() string {
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

func isGitRepo() bool {
	c := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	res, err := c.CombinedOutput()
	flag := strings.TrimSuffix(string(res), "\n")
	if err == nil && string(flag) == "true" {
		return true
	}
	fmt.Println(red, flag, reset)
	return false
}

func findReleaseBranches()([]string, error)  {
	c := exec.Command("bash", "-c", "git ls-remote --q | grep refs/heads/release/ | awk '{print $2}'")
	res, err := c.CombinedOutput()
	if err != nil {
		fmt.Println(red, err, reset)
		return nil,err
	}

	branches := string(res)
	if len(branches) == 0 {
		return nil, errors.New("没有 release 分支")
	}

	branchesArray := strings.Split(strings.Trim(branches, "\n"), "\n")
	if len(branchesArray) == 0 {
		err = errors.New("没有 release 分支")
		fmt.Println(err)
		return nil, err
	}

	r := make([]string, 0)

	for _, v := range branchesArray {
		v = strings.Trim(v, "\n")
		b := strings.Split(v, "refs/heads/")
		r = append(r,b[1])
	}

	return r, nil
}

func getTargetReleaseBranch(branches []string) string  {

	if len(branches) == 1 {
		return  branches[0]
	} else {
		fmt.Println("选择将要合入的分支")
		for i, branch := range branches {
			fmt.Println(i, branch)
		}

		var indexString string
		fmt.Scanln(&indexString)

		index, err := strconv.Atoi(indexString)

		if err != nil {
			fmt.Println("输入有误")
			return ""
		}

		if index < len(branches) {
			return branches[index]
		} else {
			fmt.Println("输入有误")
			return  ""
		}
	}
}
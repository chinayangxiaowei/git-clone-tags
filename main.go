package main

import (
	"flag"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

import "github.com/go-cmd/cmd"
import "github.com/tidwall/match"

func gitInitAndSetRemoteOrigin(gitUrl string) {
	{
		cmd := exec.Command("git", "init")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
	{
		cmd := exec.Command("git", "remote", "add", "origin", gitUrl)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}

func gitCheckTag(tag string) bool {
	c := cmd.NewCmd("git", "tag", "--list", tag)
	s := <-c.Start()
	out := s.Stdout
	if len(out) == 0 {
		return false
	}
	return out[0] == tag
}

func gitFetchTag(tag string) {
	cmd := exec.Command("git", "fetch", "origin", "tag", tag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	cmd = exec.Command("git", "checkout", "-b", "branch_"+tag, tag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func gitListEndTag(gitUrl string, endBuild bool, minMajor int, minBuild int, tagsPattern string) []string {
	var sortParam string
	//if endBuild {
	sortParam = "--sort=version:refname"
	//}
	c := cmd.NewCmd("git", "ls-remote", "--tags", sortParam, gitUrl)
	s := <-c.Start()
	out := s.Stdout
	var verMap map[string]string
	verMap = make(map[string]string)
	var lastKey, lastVer string
	var maxKeyArr []string
	comTags := len(tagsPattern) > 0
	for i := 0; i < len(out); i++ {
		namePos := strings.LastIndexByte(out[i], '/')
		if namePos == -1 {
			continue
		}
		ss := out[i][namePos+1:]
		if comTags && !match.Match(ss, tagsPattern) {
			continue
		}
		if !endBuild && minBuild == 0 {
			maxKeyArr = append(maxKeyArr, ss)
			continue
		}
		ssArr := strings.Split(ss, ".")
		if len(ssArr) == 1 {
			continue
		}
		var sKey string
		var sMajorVer string
		var iMajorVer int
		var sBuildVer string
		var iBuildVer int
		if len(ssArr) == 2 {
			sKey = ssArr[0]
			sMajorVer = ssArr[0]
			sBuildVer = ssArr[1]
		} else if len(ssArr) == 3 {
			sKey = ssArr[0] + "." + ssArr[1]
			sMajorVer = ssArr[0]
			sBuildVer = ssArr[2]

		} else if len(ssArr) == 4 {
			sKey = ssArr[0] + "." + ssArr[1] + "." + ssArr[2]
			sMajorVer = ssArr[0]
			sBuildVer = ssArr[3]
		}
		if minMajor != 0 {
			if v, err := strconv.Atoi(sMajorVer); err == nil {
				iMajorVer = v
			}
			if iMajorVer < minMajor {
				continue
			}
		}
		if v, err := strconv.Atoi(sBuildVer); err == nil {
			iBuildVer = v
		}
		if iBuildVer < minBuild {
			continue
		}
		if !endBuild {
			maxKeyArr = append(maxKeyArr, ss)
			continue
		}
		if verMap[sKey] == "" {
			verMap[sKey] = ssArr[3]
			if lastKey != "" {
				maxKeyArr = append(maxKeyArr, lastKey+"."+lastVer)
			}
		} else {
			if v, err1 := strconv.Atoi(verMap[sKey]); err1 == nil {
				if iBuildVer > v {
					verMap[sKey] = ssArr[3]
				}
			}
		}
		lastKey = sKey
		lastVer = sBuildVer
	}
	return maxKeyArr
}

func getSavePathFromUrl(gitUrl string) string {
	if u, err := url.Parse(gitUrl); err == nil {
		filename := filepath.Base(u.Path)
		ext := path.Ext(filename)
		if strings.ToLower(ext) == ".git" {
			return strings.Replace(filename, ext, "", 1)
		}
		return filename
	}
	return ""
}

func PathExists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func main() {
	var gitUrl string
	var savePath string
	var endBuild bool
	var tagsPattern string
	var showTags bool
	var minMajor int
	var minBuild int
	flag.StringVar(&gitUrl, "remote", "", "remote url (needed)")
	flag.StringVar(&savePath, "repo", "", "repository path")
	flag.StringVar(&tagsPattern, "tags", "", "tags matching string")
	flag.BoolVar(&endBuild, "end-build", false, "calculating the maximum build version")
	flag.IntVar(&minMajor, "min-major", 0, "filter minimum major version number")
	flag.IntVar(&minBuild, "min-build", 0, "filter minimum build version number")
	flag.BoolVar(&showTags, "show-tags", false, "display matching tags, but do not clone to local repo")
	flag.Parse()
	if !flag.Parsed() || gitUrl == "" {
		flag.Usage()
		os.Exit(-1)
	}
	if len(savePath) == 0 {
		savePath = getSavePathFromUrl(gitUrl)
	}
	if len(savePath) > 0 && !path.IsAbs(savePath) {
		if _path, err := filepath.Abs(savePath); err == nil {
			savePath = _path
		}
	}
	//fmt.Println(gitUrl, savePath)
	maxKeyArr := gitListEndTag(gitUrl, endBuild, minMajor, minBuild, tagsPattern)
	if len(maxKeyArr) > 0 {
		if showTags {
			for i := 0; i < len(maxKeyArr); i++ {
				fmt.Println(maxKeyArr[i])
			}
			os.Exit(0)
		}
		isExist, err := PathExists(savePath)
		if !isExist && err != nil {
			fmt.Println("Failed: ", err.Error())
			os.Exit(-1)
		}
		if !isExist {
			err := os.MkdirAll(savePath, fs.ModePerm)
			if err != nil {
				fmt.Println("Failed: ", err.Error())
				os.Exit(-1)
			}
		}
		oldPath, _ := os.Getwd()
		err = os.Chdir(savePath)
		if err != nil {
			fmt.Println("Failed: ", err.Error())
			os.Exit(-1)
		}
		isExist, err = PathExists(savePath + "/.git")
		if !isExist {
			gitInitAndSetRemoteOrigin(gitUrl)
		}
		for i := 0; i < len(maxKeyArr); i++ {
			if !gitCheckTag(maxKeyArr[i]) {
				fmt.Println(maxKeyArr[i] + " fetch")
				gitFetchTag(maxKeyArr[i])
			} else {
				fmt.Println(maxKeyArr[i] + " is exist")
			}
		}
		os.Chdir(oldPath)
	} else {
		fmt.Println("The list of remote tags is empty.")
	}
}

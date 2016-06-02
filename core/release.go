package core

import (
	"net/http"
	"errors"
	"io/ioutil"
	"regexp"
	"fmt"
	"os"
)

const apiBase = "https://data.services.jetbrains.com/products/releases?code="

var LINUX_MATCH = regexp.MustCompile("\"linux\":\\{((.|\\n)*?)\\}")
var MAC_MATCH = regexp.MustCompile("\"mac\":\\{((.|\\n)*?)\\}")
var WINDOWS_MATCH = regexp.MustCompile("\"windows\":\\{((.|\\n)*?)\\}")

var DOWNLOAD_MATCH = regexp.MustCompile("\"link\":\"(.*?)\"")
var CHECKSUM_MATCH = regexp.MustCompile("\"checksumLink\":\"(.*?)\"")
var VERSION_MATCH = regexp.MustCompile("\"version\":\"(.*?)\"")

var releaseCode = map[string]string {
	"appcode" : "AC",
	"clion" : "CL",
	"datagrip" : "DG",
	"idea" : "IIU",
	"intellij" : "IIU",
	"idea-ce" : "IIC",
	"intellij-ce" : "IIC",
	"mps" : "MPS",
	"phpstorm" : "PS",
	"pycharm" : "PCP",
	"pycharm-ce" : "PCC",
	"rubymine" : "RM",
	"webstorm" : "WS",
}

type Release struct {
	ChecksumURL string
	DownloadURL string
	Version string
}

func GetReleaseInfo(tool,dist string) (*Release,error){
	code := releaseCode[tool]
	if code == "" {
		return nil,errors.New(tool + " is an invalid toolname")
	}
	//get the JSON formatted release information
	resp, err := http.Get(apiBase+code)
	if err != nil {
		return nil,errors.New("Could not get release information, reason: " + err.Error())
	}

	//turn it into a string
	infoRaw,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,errors.New("Could not read release information response, reason: " + err.Error())
	}

	r := &Release{}
	info := string(infoRaw)

	//get version number
	version := VERSION_MATCH.FindStringSubmatch(info)
	if !(len(version) == 2 && len(version[1]) > 0) {
		return nil, errors.New("Could not find a valid version number")
	}
	r.Version = version[1]

	//get os-specific information
	var distRelease []string
	switch dist {
	case "linux":
		distRelease = LINUX_MATCH.FindStringSubmatch(info)
	case "mac":
		distRelease = MAC_MATCH.FindStringSubmatch(info)
	case "windows":
		distRelease = WINDOWS_MATCH.FindStringSubmatch(info)
	default:
		distRelease = make([]string,0)
	}

	if len(distRelease) != 3 {
		return nil, errors.New("Distribution \"" + dist + "\" unsupported for \"" + tool + "\"")
	}

	//get download link
	download := DOWNLOAD_MATCH.FindStringSubmatch(distRelease[1])
	if len(download) != 2 {
		return nil, errors.New("Could not find download link")
	}
	r.DownloadURL = download[1]

	//get checksum link
	checksum := CHECKSUM_MATCH.FindStringSubmatch(distRelease[1])
	if len(checksum) != 2 {
		return nil, errors.New("Could not find checksum link")
	}
	r.ChecksumURL = checksum[1]
	return r,nil
}

func PrintRelease(tool,dist string) {
	release,err := GetReleaseInfo(tool,dist)
	if err != nil {
		fmt.Fprintln(os.Stderr,err.Error())
		os.Exit(1)
	}
	fmt.Printf("Version: %s\n",release.Version)
	fmt.Printf("Download: %s\n",release.DownloadURL)
	fmt.Printf("Checksum: %s\n",release.ChecksumURL)
}
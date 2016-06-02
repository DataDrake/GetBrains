package core

import (
	"os"
	"path/filepath"
	"errors"
	"fmt"
	"archive/tar"
	"compress/gzip"
	"io"
	"strings"
)

func exists(path string) bool {
	_,err := os.Stat(path)
	return err == nil
}

func getInstallLocation(tool,dist,version string) (string,error) {
	path := ""
	switch dist {
	case "linux":
		path = string(filepath.Separator) + "opt"
	default:
		return path,errors.New("Distribution unsupported")
	}
	path = filepath.Join(path,tool,version)
	return path,nil
}

func InstallLinux(input,output string) error {
	raw, err := os.Open(input)
	if err == nil {
		defer raw.Close()
		g,err := gzip.NewReader(raw)
		if err == nil {
			t := tar.NewReader(g)
			err := os.MkdirAll(output,0744)
			if err == nil {
				for {
					header, err := t.Next()
					if err == io.EOF {
						break
					} else if err != nil {
						return err
					}

					name := header.Name
					n := strings.Index(name,string(filepath.Separator))
					name = name[n:len(name)]

					root := filepath.Dir(name)
					if len(root) > 0 {
						root = filepath.Join(output,root)
						if err = os.MkdirAll(root,0744); err != nil {
							return err
						}
					}

					path := filepath.Join(output, name)
					info := header.FileInfo()

					file, err:= os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
					if err != nil {
						return err
					}
					defer file.Close()
					_, err = io.Copy(file, t)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return err
}

func InstallTool(tool,dist,version string) {

	f := filepath.Join(os.TempDir(), tool)
	i, err := getInstallLocation(tool, dist,version)
	fmt.Printf("Installing %s to '%s'...\n",tool,i)
	if err == nil {
		if !exists(i) {
			switch dist {
			case "linux":
				err = InstallLinux(f,i)
			default:
				err = errors.New("Distribution not supported")
			}
		} else {
			err = errors.New("Cannot overwrite existing installation")
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr,err.Error())
		os.Exit(1)
	}
	fmt.Println("Success!")
}

func UninstallTool(tool,dist,version string) {
	i, err := getInstallLocation(tool, dist,version)
	if err == nil {
		err = os.RemoveAll(i)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr,err.Error())
		os.Exit(1)
	}
}

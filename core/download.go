package core

import (
	"os"
	"net/http"
	"io"
	"path/filepath"
	"io/ioutil"
	"errors"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func verify(tool,checksum string) error {

	t, err := os.Open(tool)
	if err != nil  {
		return err
	}
	defer t.Close()

	c, err := os.Open(checksum)
	if err != nil  {
		return err
	}
	defer c.Close()

	sum,err := ioutil.ReadAll(c)
	if err != nil {
		return errors.New("Failed to read checksum")
	}
	sum = sum[0:64]

	sha := sha256.New()
	sha.Reset()

	_,err = io.Copy(sha,t)
	if err != nil {
		return err
	}
	found := hex.EncodeToString(sha.Sum(nil))

	if found != string(sum) {
		return errors.New("File does not match checksum")
	}
	return nil
}

func getFile(path,url string) error {
	// Create the file
	out, err := os.Create(path)
	if err != nil  {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil  {
		return err
	}
	return nil
}

func DownloadCleanup(tool string) {
	f := filepath.Join(os.TempDir(),tool)
	c := f + ".sha256"
	fmt.Printf("Removing temporary file '%s'...\n",f)
	err := os.Remove(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not remove temporary file '%s', reason: %s \n",f,err.Error())
	}
	fmt.Printf("Removing temporary file '%s'...\n",c)
	err = os.Remove(c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not remove temporary file '%s', reason: %s \n",c,err.Error())
	}
}

func Download(tool string, r *Release) error {
	f := filepath.Join(os.TempDir(),tool)
	c := f + ".sha256"
	fmt.Printf("Downloading %s %s to '%s'...\n",tool,r.Version,f)
	err := getFile(f,r.DownloadURL)
	if err == nil {
		fmt.Printf("Downloading Checksum to '%s' ...\n",c)
		err = getFile(c,r.ChecksumURL)
		if err == nil {
			fmt.Println("Verifying against Checksum...")
			err = verify(f,c)
		}
	}
	return err
}

func DownloadTool(tool,dist string) *Release {
	fmt.Println("Retrieving Release Information...")
	release,err := GetReleaseInfo(tool,dist)
	if err == nil {
		err = Download(tool, release)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr,err.Error())
		os.Exit(1)
	}
	fmt.Println("Success!")
	return release
}

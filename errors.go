package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// displayHTTPErrMsg dislays an http error message,
// while displaying the master page.
func (ws *website) displayHTTPErrMsg(w http.ResponseWriter, errCode int) {

	indexBtyes := ws.getRawMasterPage()

	errText := fmt.Sprintf("Error %d - %s", errCode, http.StatusText(errCode))

	// Get the error code page
	physPath := fmt.Sprintf("%s/html/errors/gen-error.html", ws.InstallPath)
	errBtyes, err := ioutil.ReadFile(physPath)
	if err != nil {
		fmt.Println(err)
	} else {
		errBtyes = bytes.Replace(errBtyes, []byte("{{.ErrorText}}"), []byte(errText), -1)
		indexBtyes = bytes.Replace(indexBtyes, []byte("{{.ErrorPlaceHolder}}"), errBtyes, -1)
	}

	w.WriteHeader(errCode)
	w.Write(indexBtyes)
}

func (ws *website) isPathMaster(rPath string) bool {

	v := strings.Split(rPath, "/")
	masterInURL := fmt.Sprintf("/%s", v[1])

	for i := 0; i < len(ws.MaseterPages); i++ {
		if masterInURL == ws.MaseterPages[i] {
			return true
		}
	}

	return false
}

// pageNotFound tries to find the html file on disk. It will write
// a custom 404 message to the response. Note that at this point the master
// file has already been loaded, so the error message is displayed on the
// top of the master content.
func (ws *website) pageNotFound(w http.ResponseWriter, rPath string) bool {

	rPath = ws.addSuffix(rPath)

	physPath := ""

	if rPath == "/" {
		rPath = defaultHomePageName
	}

	physPath = ws.getRelPathPhysLocation(rPath)

	if !fileOrDirExists(physPath) {

		ws.displayHTTPErrMsg(w, 404)

		return true
	}

	return false
}

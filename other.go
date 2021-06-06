package main

import (
	"fmt"
	"io/ioutil"

	"strings"
)

// addSuffix adds .html to the path if not present.
func (ws *website) addSuffix(rPath string) string {

	rPath = strings.ToLower(rPath)

	rPath = strings.Replace(rPath, "/null", "/", -1)

	for i := 0; i < len(ws.MaseterPages); i++ {
		if rPath == ws.MaseterPages[i] {
			return fmt.Sprintf("%s/index.html", ws.MaseterPages[i])
		}
	}

	if !strings.HasSuffix(rPath, ".html") {
		endsWithSlash := strings.HasSuffix(rPath, "/")
		if !endsWithSlash {
			rPath = fmt.Sprintf("%s.html", rPath)
		} else {
			rPath = fmt.Sprintf("%sindex.html", rPath)
		}
	}

	return rPath
}

// getFileFromCache get the byte array of the html file, either from
// disk or from the list of byte arrays saved in memory.
func (ws *website) getFileFromCache(rPath string, physPath string) []byte {

	b := ws.Webutil.Webcache.GetItem(rPath)

	if len(b) == 0 {
		b, _ = ioutil.ReadFile(physPath)
		ws.Webutil.Webcache.AddItem(rPath, b, ws.Webutil.CacheDuration)
	}

	return b
}

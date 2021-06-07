package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"strings"
)

func (ws *website) getCurrentMasterName(masterFileName string) string {
	return fmt.Sprintf("/%s", strings.Replace(masterFileName, ".html", "", -1))
}

// setMasterFilePath maps the master page accordingly to the URL path.
func (ws *website) setMasterFilePath(r *http.Request) bool {

	masterChanged := false

	rPath := strings.ToLower(r.URL.Path)

	curMaster := ws.getCurrentMasterName(ws.MasterHTMLFileName)

	if rPath == "/" || rPath == "/index.html" {
		// Set to default.
		if ws.MasterHTMLFileName != defaultMasterPageName {
			masterChanged = true
		}
		ws.MasterHTMLFileName = defaultMasterPageName
	} else {
		// Master page has to be the first segment of the URL.
		// i.e. mymasterpage is the target in URL: http:localhost:8080/mymasterpage
		v := strings.Split(rPath, "/")
		masterInURL := fmt.Sprintf("/%s", v[1])

		for i := 0; i < len(ws.MaseterPages); i++ {
			if masterInURL == ws.MaseterPages[i] {
				if curMaster != masterInURL {
					masterChanged = true
				}
				m := strings.Replace(ws.MaseterPages[i], "/", "", -1)
				ws.MasterHTMLFileName = fmt.Sprintf("%s.html", m)
				break
			}
		}
	}

	ws.HomePageName = fmt.Sprintf("/%s", defaultHomePageName)

	return masterChanged
}

// serveRoot processing all requests for your html content.
// Note that request to your assets (js, css, img files)
// are done by ws.Webutil.ServeStaticFile().
func (ws *website) serveRoot(w http.ResponseWriter, r *http.Request) {

	serveFileWithNoMasterPage := false

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method == "HEAD" {
		w.WriteHeader(200)
		// nothing to return
		// You could set your defautlt response headers here.
		return
	}

	if !ws.validateRequest(w, r) {
		return
	}

	if ws.Config.MaintenanceWindowOn {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		p := fmt.Sprintf("%s/html/maint-window.html", ws.AppDataPath)
		b, _ := ioutil.ReadFile(p)
		b = bytes.ReplaceAll(b, []byte("{{.HostName}}"), []byte(ws.Config.HostName))
		b = bytes.Replace(b, []byte("{{.ThisYear}}"), []byte(fmt.Sprintf("%d", time.Now().Year())), -1)
		b = ws.Webutil.RemoveCommentsFromByBiteArry(b, "{{.COMMENT ", "}}")
		w.WriteHeader(200)
		w.Write(b)
		return
	}

	// The end user can type just the name of the page or name+.html.
	// i.e. ../contacts or ../contacts.html.
	rPath := ws.addSuffix(r.URL.Path)

	// Outsite means that the whole page will have to be
	// rereshed using a single html content vs a master page.
	// To make the distinction, all such paths are saved in
	// ../pages/serve-file. And the link begings with ../x/.
	// e.g. http://localhost:8000/x/myoutsidepage.html.
	// See if the page matche server-file
	physPath := fmt.Sprintf("%s/html/pages/serve-file%s", ws.InstallPath, rPath)
	if fileOrDirExists(physPath) {
		serveFileWithNoMasterPage = true
	}

	// Prevent the call to get to these directories.
	doNotServerPath := []string{"/appdata", "/html", "/errors"}
	for i := 0; i < len(doNotServerPath); i++ {
		if strings.HasPrefix(rPath, doNotServerPath[i]) {
			ws.displayHTTPErrMsg(w, 404)
			return
		}
	}

	// sitemap.xml and robots.txt
	if rPath == "/sitemap.xml" || rPath == "/robots.txt" {
		physPath := fmt.Sprintf("%s%s", ws.InstallPath, rPath)
		bXML, err := ioutil.ReadFile(physPath)
		if err == nil {
			w.Write(bXML)
			return
		}
	}

	if ws.setMasterFilePath(r) {
		// Reload the master page
		ws.serveMasterPage(w, r)
		return
	}

	if ws.pageNotFound(w, r.URL.Path) {
		// Response is already written by ws.processPageNotFound().
		return
	}

	// Server the file as-is. This is not part of the master file,
	// althought it could have its own master.
	if serveFileWithNoMasterPage {
		b, err := ioutil.ReadFile(physPath)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		b = bytes.Replace(b, []byte("{{.RootFullURL}}"), []byte(ws.RootFullURL), -1)
		b = bytes.Replace(b, []byte("{{.ThisYear}}"), []byte(fmt.Sprintf("%d", time.Now().Year())), -1)

		w.Write(b)

		return
	}

	// All content requests have this header: X-Page-Key.
	//
	//  1. The master is loaded first.
	//  2. Requests for contents go in the placeholder inside in a quebunes:div tag:
	//     <quebunes:div class="master" id="{{.MasterContentID}}">
	//     </quebunes:div>
	//
	key := r.Header.Get("X-Page-Key")

	if key != ws.PageKey {
		// First time the page is loaded or the page is refreshed by the user.
		ws.serveMasterPage(w, r)
		return
	}

	// This will go into the innerHTML of the <quebunes:div> tag.
	b := ws.setSinglePage(r.URL.Path)

	w.Write(b)
}

// getRawMasterPage gets the bare master file. The master html file contains
// the minimum html body that is to disply on your site. Note that the master
// page loads once, and navigation after that occurs is done via XMLHttpRequest
// without re-rendering any part of the page. If the end-user refreshes the page
// the master html file will load again + the last navigated html subject --
// see setSinglePage().
func (ws *website) getRawMasterPage() []byte {

	// Get the master file as-is.
	physPath := fmt.Sprintf("%s/html/pages/master/%s", ws.InstallPath, ws.MasterHTMLFileName)

	// Read the file directly or from the cache.
	// Always read from disk.
	// If you want to cache the master files, you'll have to cache all of them.
	masterBtyes, err := ioutil.ReadFile(physPath)
	if err != nil {
		fmt.Println("getRawMasterPage()", err)
	}

	// Apply your variables here (resembles the template operation)
	if len(ws.GetContentJavscript) == 0 {
		// Cache this here for performance; and more importantly so that templates are to be update
		// and implemented after some testing... so the server has to be restareted.
		tPath := fmt.Sprintf("%s/html/templates/getcontent.js", ws.InstallPath)
		ws.GetContentJavscript, _ = ioutil.ReadFile(tPath)
	}

	ws.GetContentJavscript = bytes.Replace(ws.GetContentJavscript, []byte("{{.RootFullURL}}"), []byte(ws.RootFullURL), -1)
	ws.GetContentJavscript = bytes.Replace(ws.GetContentJavscript, []byte("{{.PageKey}}"), []byte(ws.PageKey), -1)
	ws.GetContentJavscript = bytes.Replace(ws.GetContentJavscript, []byte("{{.MasterContentID}}"), []byte(ws.MasterContentID), -1)

	masterBtyes = bytes.Replace(masterBtyes, []byte("{{.ThisYear}}"), []byte(fmt.Sprintf("%d", time.Now().Year())), -1)
	masterBtyes = bytes.Replace(masterBtyes, []byte("{{.MasterContentID}}"), []byte(ws.MasterContentID), -1)
	masterBtyes = bytes.Replace(masterBtyes, []byte("{{.RootFullURL}}"), []byte(ws.RootFullURL), -1)
	masterBtyes = bytes.Replace(masterBtyes, []byte("{{.RawQueryPlaceHolder}}"), []byte(ws.RawQueryPlaceHolder), -1)
	masterBtyes = bytes.Replace(masterBtyes, []byte("{{.JSGetContent}}"), []byte(ws.GetContentJavscript), -1)
	masterBtyes = bytes.Replace(masterBtyes, []byte("{{.PageKey}}"), []byte(ws.PageKey), -1)

	return masterBtyes
}

func (ws *website) getRelPathPhysLocation(rPath string) string {

	// Oustide of master page?
	physPath := fmt.Sprintf("%s/html/pages/serve-file%s", ws.InstallPath, rPath)
	if fileOrDirExists(physPath) {
		return physPath
	}

	isMasterDefault := ws.isPathMaster(rPath)

	if !isMasterDefault {
		curMaster := ws.getCurrentMasterName(defaultMasterPageName)
		physPath = fmt.Sprintf("%s/html/pages%s%s", ws.InstallPath, curMaster, rPath)
	} else {
		physPath = fmt.Sprintf("%s/html/pages%s", ws.InstallPath, rPath)
	}

	return physPath
}

// setSinglePage returns the html text in form of []byte.
// This is the content that is display dynamically in the
// master page, without having to refresh or re-render the page.
func (ws *website) setSinglePage(rPath string) []byte {

	var bPage []byte

	rPath = strings.ToLower(rPath)
	rPath = ws.addSuffix(rPath)

	if rPath == "" {
		return bPage
	}

	physPath := ws.getRelPathPhysLocation(rPath)

	if !fileOrDirExists(physPath) {

		// Try adding home page
		physPath = fmt.Sprintf("%s%s/%s", ws.InstallPath, rPath, ws.HomePageName)
		if !fileOrDirExists(physPath) {
			return bPage
		}
	}

	//bPage = ws.getFileFromCache(rPath, physPath)
	bPage, _ = ioutil.ReadFile(physPath)

	pageName := rPath[1:] // minus the first /

	if len(ws.TopJavascript) == 0 {
		// Cache this here for performance; and more importantly so that templates are to be update
		// and implemented after some testing... so the server has to be restareted.
		tPath := fmt.Sprintf("%s/html/templates/top_js.html", ws.InstallPath)
		ws.TopJavascript, _ = ioutil.ReadFile(tPath)
	}

	topJS := ws.TopJavascript
	topJS = bytes.Replace(topJS, []byte("{{.RootFullURL}}"), []byte(ws.RootFullURL), -1)

	topJS = bytes.Replace(topJS, []byte("{{.MasterContentID}}"), []byte(ws.MasterContentID), -1)
	topJS = bytes.Replace(topJS, []byte("{{.PageName}}"), []byte(pageName), -1)
	bPage = bytes.Replace(bPage, []byte("{{.TopJS}}"), topJS, -1)

	bPage = bytes.Replace(bPage, []byte("{{.ThisYear}}"), []byte(fmt.Sprintf("%d", time.Now().Year())), -1)
	bPage = bytes.Replace(bPage, []byte("{{.RootFullURL}}"), []byte(ws.RootFullURL), -1)
	// Set the erorr tag to empty.
	var empty []byte
	bPage = bytes.Replace(bPage, []byte("{{.ErrorPlaceHolder}}"), empty, -1)

	return bPage
}

// serveMasterPage loads the homepage + the target content page.
func (ws *website) serveMasterPage(w http.ResponseWriter, r *http.Request) {

	masterBtyes := ws.getRawMasterPage()

	// Set the erorr tag to empty.
	var empty []byte
	masterBtyes = bytes.Replace(masterBtyes, []byte("{{.ErrorPlaceHolder}}"), empty, -1)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Encoding", "compress")
	w.Write(masterBtyes)
}

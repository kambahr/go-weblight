package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kambahr/go-webcache"
	"github.com/kambahr/go-webconfig"

	"time"

	"github.com/kambahr/go-webutil"
)

func main() {

	var portNoArg uint

	mWebsite := new(website)

	flag.UintVar(&portNoArg, "portno", 8080, "tcp/ip port no to listen - defaults to port 8080")
	flag.Parse()

	// Current directory is the default path. The EXE could run
	// in one folder and the path to the web files (html, js,..)
	// could be entirely a different path (i.e. on a differe volume).
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	mWebsite.start(portNoArg, dir)
}

// func (ws *website) foo(ipAddr string) {
// 	TODO:
// }
func (ws *website) connState(conn net.Conn, connState http.ConnState) {

	// You can block target IP addresses here by closing connection of
	// the offenders or the ones that are not be in your authorized group.

	// This is the <IP address>:<identifier> e.g. 127.0.0.1:35692
	// The numbers after the colon is a unique id given to each requested page.
	//    For example:
	//    http://localhost:8000/mycss.css........127.0.0.1:35692
	//    http://localhost:8000/myjs.js..........127.0.0.1:35693
	//
	// If you use XMLHttpRequest() to get a response back. The identifier will remain
	// the same even though you may get a different content (different html code).
	//    For example:
	//    http://localhost:8000/page1.html........127.0.0.1:35692
	//    http://localhost:8000/page2.html........127.0.0.1:35692
	//
	// For the most part, you can ignore the identifier, however, you can make use
	// of it -- to build on security, content/performance-smart concepts to enhance
	// your site.

	// If you want to do any processing here, you'd have to use the go statement
	// so that the processing is done outside of this function (in a speparate
	// thread).
	//go ws.foo(conn.RemoteAddr().String())

	maxConnIODeadLine := 20

	// Some extra for over timeout.
	maxConnIOWriteDeadLine := 30

	// These help reduce the too many open files errors.
	conn.SetDeadline(time.Now().Add(time.Duration(maxConnIODeadLine) * time.Second))
	conn.SetReadDeadline(time.Now().Add(time.Duration(maxConnIOWriteDeadLine) * time.Second))
	conn.SetWriteDeadline(time.Now().Add(time.Duration(maxConnIOWriteDeadLine) * time.Second))

	// Uncomment this if you want to do anything on active, idle, and closed states.

	// csStr := fmt.Sprintf("%v", connState)
	// if csStr == "new" {
	// } else if csStr == "active" {
	// } else if csStr == "idle" {
	// } else if csStr == "closed" {
	// }
}

func parseCofigLine(line string, key string) string {
	return strings.Trim(strings.Replace(line, key, "", 1), " ")
}

// start initializes your website.
func (ws *website) start(portNo uint, appPath string) {

	ws.InstallPath = appPath

	// Create the appdata if it does not exist
	ws.AppDataPath = fmt.Sprintf("%s/appdata", ws.InstallPath)
	if !fileOrDirExists(ws.AppDataPath) {
		os.Mkdir(ws.AppDataPath, os.ModePerm)
	}

	// Initializes an instance of webconfig. e.Config is refreshed
	// automatically. To update values:
	// e.Config.UpdateConfigValue(<key>, <value>).
	// See https://github.com/kambahr/go-webconfig.
	ws.Config = webconfig.NewWebConfig(ws.InstallPath)

	if portNo > 0 {
		ws.Config.UpdateConfigValue("", "portno", fmt.Sprintf("%d", portNo))
		ws.Config.Refresh()
	}

	ws.RootFullURL = fmt.Sprintf("%s://%s:%d", ws.Config.Proto, ws.Config.HostName, ws.Config.PortNo)

	// Keep this constant; it's just a placeholder and does not
	// have any security attributes.
	ws.PageKey = "a7d7wka26bf0f508f4b365f94b25bc1f35add97d529aba1f0"
	ws.RawQueryPlaceHolder = fmt.Sprintf("%s%s", getMd5String(time.Now().String()), getRandCharacter())

	// This is good to change everytime the webserver starts
	ws.MasterContentID = getRandCharacter()
	ws.MasterHTMLFileName = defaultMasterPageName
	ws.HomePageName = defaultHomePageName

	// Configure the master pages' location.
	// Master pages are placed in folders under Pages. Antying unter the root
	// is under the / path (i.e. mypage.html => http://localhost:8080/mypage.html).
	t := fmt.Sprintf("%s/html/pages", ws.InstallPath)
	dirs, _ := getDirectories(t)
	for i := 0; i < len(dirs); i++ {
		if strings.HasPrefix(dirs[i], "/pages") {
			continue
		}
		p := strings.Replace(dirs[i], t, "", 1)
		if p != "" && p != "/serve-file" && p != "/master" {
			ws.MaseterPages = append(ws.MaseterPages, p)
		}
	}

	ws.Webcache = webcache.NewWebCache(5 * time.Minute)

	ws.Webutil = webutil.NewHTTP(ws.InstallPath, 5*time.Minute)

	// This covers all supporting files like js, css, img,...
	http.HandleFunc("/assets/", ws.Webutil.ServeStaticFile)

	// And this covers all of the logic for your site: static/rendered html,
	// background, redirects, etc...
	http.HandleFunc("/", ws.serveRoot)

	var svr http.Server
	{
		svr = http.Server{
			Addr:           fmt.Sprintf(":%d", portNo),
			MaxHeaderBytes: 20480,

			// This prevents http2 calls by the client browser
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),

			// No checks on the client cert
			TLSConfig: &tls.Config{
				ClientAuth: tls.NoClientCert,
			},

			// These should reflect what you expect from your site
			// also see comments in ws.ConnStatex().
			ReadTimeout:       20 * time.Second,
			WriteTimeout:      20 * time.Second,
			IdleTimeout:       60 * time.Second,
			ReadHeaderTimeout: 10 * time.Second,

			ConnState: ws.connState,
		}
	}

	fmt.Println("Listening to port:", portNo, "for host:", ws.Config.HostName)
	if ws.Config.Proto == "HTTP" {
		log.Fatal(svr.ListenAndServe())
	} else {
		if fileOrDirExists(ws.Config.TLS.CertFilePath) &&
			fileOrDirExists(ws.Config.TLS.KeyFilePath) {

			// both files must be in PEM format.
			//   Note that there is no separate arg to apply the
			//   chaine certficate. The one cert file can include the
			//   chain (if any) along with the root cert.

			// Assume that they are valid and run
			log.Fatal(svr.ListenAndServeTLS(ws.Config.TLS.CertFilePath, ws.Config.TLS.KeyFilePath))
		}
	}
}

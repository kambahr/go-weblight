package main

import (
	"github.com/kambahr/go-webcache"
	"github.com/kambahr/go-webconfig"

	"github.com/kambahr/go-webutil"
)

const (
	defaultMasterPageName = "default.html"
	defaultHomePageName   = "index.html"
)

type tlsFiles struct {
	CertFilePath string
	KeyFilePath  string
}

// website holds all the global vars and types needed to run a website.
type website struct {

	// MaseterPages is a list of directory structure of the master pages.
	MaseterPages []string

	// MasterHTMLFileName is the file name of the master (container) file that
	// that contains all dynamic content. For a simple setting this is can be
	// called master.html.
	// You can load different master files, while your site is running.
	// This could be done to replace the entire logic, or to chnage
	// the appearance of your site -- or both.
	MasterHTMLFileName string

	// HomePageName is a default landing page. It's usally named index.html.
	HomePageName string

	//PageKey is used to see if the master page has already been loaded.
	PageKey string

	//RawQueryPlaceHolder is a place-holder query string to submit to the same page.
	RawQueryPlaceHolder string

	// Webutil helps in serving static files that to be served by default.
	// Theses files are served as-is from disk.
	Webutil *webutil.HTTP

	Config *webconfig.Config

	// Webcache allows you to a selected number of files.
	Webcache *webcache.Cache

	// Install path is the root dir to where the web files
	// (html, js, css,...) are loacted at.
	InstallPath string

	AppDataPath string

	// RootFullURL is used to place fully qualified urls in html pages
	// to avoid confilict that occur from relative url paths.
	RootFullURL string

	// TopJavascript is a javascript block that refreshes the page,
	// if the master page has not been loaded.
	TopJavascript []byte

	// GetContentJavscript is the javascript block that handles
	// the dynamic content and is added to each master file.
	GetContentJavscript []byte

	// MasterContentID is the DOM id of the <quenubes:div> tag that hold the sub-pages.
	// Note that javascript will run inside <quenubes:div>, while
	// it would not run inside a <div> tag.
	MasterContentID string
}

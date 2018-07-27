package main

import (
	"io/ioutil"
	"path"
	"strings"

	"github.com/gramework/gramework"
	"github.com/spf13/pflag"
)

var (
	indexPath   = pflag.StringP("indexpath", "i", "./index.html", "index file path")
	staticPath  = pflag.StringP("staticpath", "s", "./", "static directory path")
	staticRoute = pflag.StringP("staticroute", "r", "/static", "static route")
	bind        = pflag.StringP("bind", "b", ":80", "port to listen")
	cache       = pflag.BoolP("cache", "c", false, "enable cache")
	tls         = pflag.BoolP("tls", "t", false, "enable TLS via letsencrypt")
	dev         = pflag.BoolP(
		"dev",
		"d",
		false,
		"enables the dev mode for TLS cert generation,\nwhich allows you to use self-signed certs on localhost",
	)
)

func serveDir(app *gramework.App, sp, sr string, slashCnt int) {
	h := app.ServeDirNoCacheCustom(sp, slashCnt, true, false, []string{"index.html"})
	if *cache {
		h = app.ServeDirCustom(sp, slashCnt, true, false, []string{"index.html"})
	}

	app.GET(path.Join(sr, "/*any"), h)
}

func regStaticHandlers(app *gramework.App, sp string, slashCnt int) {
	files, err := ioutil.ReadDir(sp)
	if err != nil {
		app.Logger.WithError(err).WithField("path", sp).Fatal("static directory does not exist")
	}

	for _, file := range files {
		p := path.Clean(path.Join(sp, file.Name()))
		r := "/" + strings.TrimLeft(file.Name(), "/")

		if strings.Contains(file.Name(), ".fasthttp.gz") || path.Clean(*indexPath) == p {
			continue
		}

		if file.IsDir() {
			serveDir(app, p, r, slashCnt)
		} else {
			app.ServeFile(r, p)
		}
	}
}

func main() {
	pflag.Parse()

	gramework.DisableFlags()
	app := gramework.New()
	app.TLSEmails = []string{
		"3s@gramework.win",
	}

	app.SPAIndex(*indexPath)

	sp := path.Clean(*staticPath)
	sr := path.Clean(*staticRoute)
	slashCnt := strings.Count(sr, "/")
	if len(*staticRoute) == 0 {
		sr = "/"
	}

	if sr == "/" {
		regStaticHandlers(app, sp, slashCnt)
	} else {
		serveDir(app, sp, sr, slashCnt)
	}

	if *tls {
		if *dev {
			app.ListenAndServeAllDev(*bind)
			return
		}
		app.ListenAndServeAll()
		return
	}
	app.ListenAndServe(*bind)
}

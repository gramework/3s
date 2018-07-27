package main

import (
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
)

func main() {
	pflag.Parse()

	gramework.DisableFlags()
	app := gramework.New()

	app.SPAIndex(*indexPath)

	slashCnt := strings.Count(*staticRoute, "/")
	h := app.ServeDirNoCacheCustom(*staticPath, slashCnt, true, false, nil)
	if *cache {
		h = app.ServeDirCustom(*staticPath, slashCnt, true, false, nil)
	}
	app.GET(*staticRoute+"/*any", h)

	app.ListenAndServe(*bind)
}

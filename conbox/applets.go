package main

import (
	"github.com/udhos/conbox/applets/cat"
	"github.com/udhos/conbox/applets/echo"
	"github.com/udhos/conbox/applets/ls"
	"github.com/udhos/conbox/applets/mkdir"
	"github.com/udhos/conbox/applets/printenv"
	"github.com/udhos/conbox/applets/ps"
	"github.com/udhos/conbox/applets/pwd"
	"github.com/udhos/conbox/applets/rm"
	"github.com/udhos/conbox/applets/rmdir"
	"github.com/udhos/conbox/applets/shell"
	"github.com/udhos/conbox/applets/touch"
	"github.com/udhos/conbox/applets/which"
	"github.com/udhos/conbox/common"
)

func loadApplets() map[string]common.AppletFunc {
	tab := map[string]common.AppletFunc{
		"cat":      cat.Run,
		"echo":     echo.Run,
		"ls":       ls.Run,
		"mkdir":    mkdir.Run,
		"printenv": printenv.Run,
		"pwd":      pwd.Run,
		"ps":       ps.Run,
		"rm":       rm.Run,
		"rmdir":    rmdir.Run,
		"shell":    shell.Run,
		"touch":    touch.Run,
		"which":    which.Run,
	}
	return tab
}

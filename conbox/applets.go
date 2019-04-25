package main

import (
	"github.com/udhos/conbox/applets/cat"
	"github.com/udhos/conbox/applets/echo"
	"github.com/udhos/conbox/applets/ls"
	"github.com/udhos/conbox/applets/pwd"
	"github.com/udhos/conbox/applets/rm"
	"github.com/udhos/conbox/applets/shell"
	"github.com/udhos/conbox/common"
)

func loadApplets() map[string]common.AppletFunc {
	tab := map[string]common.AppletFunc{
		"cat":   cat.Run,
		"echo":  echo.Run,
		"ls":    ls.Run,
		"pwd":   pwd.Run,
		"rm":    rm.Run,
		"shell": shell.Run,
	}
	return tab
}

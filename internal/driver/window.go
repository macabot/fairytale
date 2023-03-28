package driver

import (
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/driver/js"
)

// TODO refactor such that driver is chosen by client?
var Window hypp.Window = js.Driver{}.Window()

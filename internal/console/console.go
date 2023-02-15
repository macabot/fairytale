package console

import "syscall/js"

func Debug(args ...any) {
	js.Global().Get("console").Call("debug", args...)
}

func Warn(args ...any) {
	js.Global().Get("console").Call("warn", args...)
}

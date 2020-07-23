package logger

import (
	"fmt"
	"github.com/kataras/golog"
	"github.com/kataras/pio"
	"strings"
)

func Hijacker(ctx *pio.Ctx) {
	l, ok := ctx.Value.(*golog.Log)
	if !ok {
		ctx.Next()
		return
	}

	w := ctx.Printer

	var logOut []string

	if l.Level != golog.DisableLevel {
		if level, ok := golog.Levels[l.Level]; ok {
			if pio.SupportColors(w.Output) {
				logOut = append(logOut, pio.Rich(level.Title, level.ColorCode, level.Style...))
			} else {
				logOut = append(logOut, level.Title)
			}
		}
	}

	if t := l.FormatTime(); t != "" {
		logOut = append(logOut, t)
	}

	if prefix := l.Logger.Prefix; len(prefix) > 0 {
		// logOut = append(logOut, prefix)
		l.Message = prefix + l.Message
	}
	if l.Logger.NewLine {
		_, _ = fmt.Fprintln(w, strings.Join(logOut, " "), l.Message)
	} else {
		_, _ = fmt.Fprint(w, strings.Join(logOut, " "), l.Message)
	}
	ctx.Store(nil, pio.ErrHandled)
}

func init() {

}

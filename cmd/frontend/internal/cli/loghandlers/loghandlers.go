package loghandlers

import (
	"gopkg.in/inconshreveable/log15.v2"
	"strings"
	"time"
)

// Trace returns  a filter for the given traces that run longer than threshold
func Trace(types []string, threshold time.Duration) func(record *log15.Record) bool {
	all := false
	valid := map[string]bool{}
	for _, t := range types {
		valid[t] = true
		if t == "all" {
			all = true
		}
	}
	return func(r *log15.Record) bool {
		if r.Lvl != log15.LvlDebug {
			return true
		}
		if !strings.HasPrefix(r.Msg, "TRACE ") {
			return true
		}
		if !all && !valid[r.Msg[6:]] {
			return false
		}

		for i := 1; i < len(r.Ctx); i += 2 {
			if r.Ctx[i-1] != "duration" {
				continue
			}
			d, ok := r.Ctx[i].(time.Duration)
			return !ok || d >= threshold
		}
		return true
	}
}

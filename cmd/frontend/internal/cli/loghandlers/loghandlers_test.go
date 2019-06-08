package loghandlers

import (
	"gopkg.in/inconshreveable/log15.v2"
	"testing"
	"time"
)

var traces = []log15.Record{
	mkRecord(log15.LvlDebug, "TRACE backend", "rpc", "RepoTree.Get", "duration", time.Second),
	mkRecord(log15.LvlDebug, "TRACE HTTP", "routename", "repo.resolve", "duration", time.Second/3),
	mkRecord(log15.LvlDebug, "TRACE HTTP", "routename", "repo.resolve", "duration", 2*time.Second),
}

func TestTrace_All(t *testing.T) {
	f := Trace([]string{"all"}, 0)
	for _, r := range traces {
		if !f(&r) {
			t.Errorf("Should allow %v", r)
		}
	}
}

func TestTrace_None(t *testing.T) {
	f := Trace([]string{}, 0)
	for _, r := range traces {
		if f(&r) {
			t.Errorf("Should filter %v", r)
		}
	}
}

func TestTrace_Specific(t *testing.T) {
	f := Trace([]string{"HTTP"}, 0)
	for _, r := range traces {
		keep := r.Msg == "TRACE HTTP"
		if f(&r) == keep {
			continue
		} else if keep {
			t.Errorf("Should keep %v", r)
		} else {
			t.Errorf("Should filter %v", r)
		}
	}
}

func mkRecord(lvl log15.Lvl, msg string, ctx ...interface{}) log15.Record {
	return log15.Record{
		Lvl: lvl,
		Msg: msg,
		Ctx: ctx,
	}
}

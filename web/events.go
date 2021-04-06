package web

import (
	"context"
	"fmt"
	"localhost/htmltoebook/types"

	"github.com/jfyne/live"
)

func onLogMsg(ctx context.Context, s *live.Socket, p map[string]interface{}) (interface{}, error) {
	m := newModel(s)
	msgparam, ok := p["logmsg"]
	if !ok {
		return m, fmt.Errorf("log not provided %v", p)
	}
	logmsg, ok := msgparam.(types.LogMsg)
	if !ok {
		return m, fmt.Errorf("log conversion error %v", p)
	}
	m.LogMsgs = []types.LogMsg{logmsg}
	return m, nil
}

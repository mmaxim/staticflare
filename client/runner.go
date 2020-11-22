package client

import (
	"time"

	"mmaxim.org/staticflare/common"
)

type Runner struct {
	*common.DebugLabeler
	remoteIPSource RemoteIPSource
}

func NewRunner(remoteIPSource RemoteIPSource) *Runner {
	return &Runner{
		DebugLabeler:   common.NewDebugLabeler("Runner"),
		remoteIPSource: remoteIPSource,
	}
}

func (r *Runner) Run() {
	for {
		ip, err := r.remoteIPSource.GetRemoteIP()
		if err != nil {
			r.Debug("Run: failed to get IP: %s", err)
		} else {

			r.Debug("Run: ip: %s", ip)
		}
		time.Sleep(time.Second)
	}
}

package client

import (
	"encoding/json"

	"mmaxim.org/staticflare/common"
)

type StaticFlaredHandler struct {
}

func NewStaticFlaredHandler() *StaticFlaredHandler {
	return &StaticFlaredHandler{}
}

func (h *StaticFlaredHandler) ParseRemoteIP(resp string) (res string, err error) {
	var info common.InfoResponse
	if err := json.Unmarshal([]byte(resp), &info); err != nil {
		return res, err
	}
	return info.RemoteIP, nil
}

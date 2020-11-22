package common

import "log"

type DebugLabeler struct {
	name string
}

func NewDebugLabeler(name string) *DebugLabeler {
	return &DebugLabeler{
		name: name,
	}
}

func (l *DebugLabeler) Debug(line string, args ...interface{}) {
	log.Printf(l.name+": "+line+"\n", args...)
}

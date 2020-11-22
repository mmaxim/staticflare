package common

import stathat "github.com/stathat/go"

type StatsProvider interface {
	Count(name string, count int)
	CountOne(name string)
	Value(name string, value float64)
	SetPrefix(prefix string) StatsProvider
}

type DummyStatsProvider struct {
}

func NewDummyStatsProvider() DummyStatsProvider {
	return DummyStatsProvider{}
}

func (d DummyStatsProvider) Count(name string, count int)          {}
func (d DummyStatsProvider) CountOne(name string)                  {}
func (d DummyStatsProvider) Value(name string, value float64)      {}
func (d DummyStatsProvider) SetPrefix(prefix string) StatsProvider { return d }

type StathatStatsProvider struct {
	*DebugLabeler
	prefix string
	ezkey  string
}

func NewStathatStatsProvider(name, ezkey string) *StathatStatsProvider {
	return &StathatStatsProvider{
		DebugLabeler: NewDebugLabeler("StathatStatsProvider"),
		prefix:       name,
		ezkey:        ezkey,
	}
}

func (s *StathatStatsProvider) fullName(name string) string {
	return s.prefix + " - " + name
}

func (s *StathatStatsProvider) SetPrefix(prefix string) StatsProvider {
	return &StathatStatsProvider{
		ezkey:  s.ezkey,
		prefix: s.prefix + " - " + prefix,
	}
}

func (s *StathatStatsProvider) Count(name string, count int) {
	if err := stathat.PostEZCount(s.fullName(name), s.ezkey, count); err != nil {
		s.Debug("Count: failed to post: %s", err)
	}
}

func (s *StathatStatsProvider) CountOne(name string) {
	if err := stathat.PostEZCountOne(s.fullName(name), s.ezkey); err != nil {
		s.Debug("CountOne: failed to post: %s", err)
	}
}

func (s *StathatStatsProvider) Value(name string, value float64) {
	if err := stathat.PostEZValue(s.fullName(name), s.ezkey, value); err != nil {
		s.Debug("Value: failed to post: %s", err)
	}
}

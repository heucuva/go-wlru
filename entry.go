package wlru

type Entry struct {
	Key         interface{}
	Value       interface{}
	IsPermanent bool
	IsExpired   bool
}

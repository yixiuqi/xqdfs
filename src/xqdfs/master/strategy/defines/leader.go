package defines

type Leader interface {
	IsLeader() bool
	LeaderId() string
	Stop()
}

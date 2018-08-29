package process

type Replication interface {
	Process(task map[int32]*ReplicationTask)
}

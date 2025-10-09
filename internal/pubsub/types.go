package pubsub

// an enum to represent "durable" or "transient"
type SimpleQueueType int

const (
	SimpleQueueTypeDurable SimpleQueueType = iota
	SimpleQueueTypeTransient
)

type AckType int

const (
	Ack AckType = iota
	NackRequeue
	NackDiscard
)

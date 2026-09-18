package outbox_relay

type eventStatus string

const (
	eventStatusPending    eventStatus = "pending"
	eventStatusProcessing eventStatus = "processing"
	eventStatusRetry      eventStatus = "retry"
	eventStatusSucceeded  eventStatus = "succeeded"
	eventStatusDead       eventStatus = "dead"
)
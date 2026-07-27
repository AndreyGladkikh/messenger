package event

type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusRetry      Status = "retry"
	StatusSucceeded  Status = "succeeded"
	StatusDead       Status = "dead"
)

package queue

import "context"

type Job struct {
	Kind      string `json:"kind"`
	AccountID uint64 `json:"account_id"`
}

type Producer interface {
	PushInitialSyncJob(ctx context.Context, accountID uint64) (jobID uint64, err error)
	Close() error
}

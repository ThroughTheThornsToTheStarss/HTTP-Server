package queue

import "context"

const (
	JobKindInitialSync   = "initial_sync"
	JobKindWebhookUpsert = "webhook_upsert"
	JobKindWebhookDelete = "webhook_delete"
)

type Job struct {
	Kind       string  `json:"kind"`
	AccountID  uint64  `json:"account_id"`
	ContactIDs []int64 `json:"contact_ids,omitempty"`
	Action     string  `json:"action"`
}

type Producer interface {
	PushInitialSyncJob(ctx context.Context, accountID uint64) (jobID uint64, err error)
	PushWebhookUpsertJob(ctx context.Context, accountID uint64, contactIDs []int64) (uint64, error)
	PushWebhookDeleteJob(ctx context.Context, accountID uint64, contactIDs []int64) (uint64, error)
	Close() error
}



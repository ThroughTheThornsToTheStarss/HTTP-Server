package beanstalk

import (
	"context"
	"encoding/json"
	"errors"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/queue"
)

func (p *Producer) put(ctx context.Context, job queue.Job) (uint64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	if job.AccountID == 0 {
		return 0, errors.New("account_id must be > 0")
	}

	if job.Kind == "" {
		return 0, errors.New("job kind is required")
	}

	body, err := json.Marshal(job)
	if err != nil {
		return 0, err
	}

	return p.tube.Put(body, defaultPriority, defaultDelay, defaultTTR)
}

func (p *Producer) Close() error {
	if p == nil || p.conn == nil {
		return nil
	}
	return p.conn.Close()
}

func (p *Producer) PushInitialSyncJob(ctx context.Context, accountID uint64) (uint64, error) {
	return p.put(ctx, queue.Job{
		Kind:      queue.JobKindInitialSync,
		AccountID: accountID,
	})
}

func (p *Producer) PushWebhookUpsertJob(ctx context.Context, accountID uint64, contactIDs []int64) (uint64, error) {
	return p.put(ctx, queue.Job{
		Kind:       queue.JobKindWebhookUpsert, 
		AccountID:  accountID,
		ContactIDs: contactIDs,
	})
}

func (p *Producer) PushWebhookDeleteJob(ctx context.Context, accountID uint64, contactIDs []int64) (uint64, error) {
	return p.put(ctx, queue.Job{
		Kind:       queue.JobKindWebhookDelete, 
		AccountID:  accountID,
		ContactIDs: contactIDs,
	})
}

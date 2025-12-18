package beanstalk

import (
	"context"
	"encoding/json"
	"errors"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/queue"
)

func (p *Producer) PushInitialSyncJob(ctx context.Context, accountID uint64) (uint64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	if accountID == 0 {
		return 0, errors.New("account_id must be > 0")
	}

	job := queue.Job{
		Kind:      "contacts_initial_sync",
		AccountID: accountID,
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


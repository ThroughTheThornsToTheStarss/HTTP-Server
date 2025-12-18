package usecase

import (
	"context"
	"errors"
	"log"
	"strings"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/domain"
	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/queue"
	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/repo"
)

type WebhookContactsUsecase struct {
	repo     repo.Repository
	producer queue.Producer
}

func NewWebhookContactsUsecase(r repo.Repository, p queue.Producer) *WebhookContactsUsecase {
	return &WebhookContactsUsecase{repo: r, producer: p}
}

type WebhookContactEvent struct {
	AmoID   int64
	Name    string
	Email   *string
	Deleted bool
}

func (u *WebhookContactsUsecase) Handle(ctx context.Context, accountID uint64, events []WebhookContactEvent) (uint64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	if accountID == 0 {
		return 0, errors.New("account_id must be > 0")
	}

	acc, err := u.repo.GetAccountByID(accountID)
	if err != nil || acc == nil || !acc.IsActive {
		return 0, nil
	}

	if u.producer == nil {
		return 0, errors.New("queue producer is nil")
	}

	upsertIDs := make([]int64, 0, len(events))
	deleteIDs := make([]int64, 0, len(events))

	for _, ev := range events {
		if ev.AmoID <= 0 {
			continue
		}

		var email *string
		if ev.Email != nil {
			v := strings.TrimSpace(*ev.Email)
			if v != "" {
				email = &v
			}
		}

		status := domain.ContactStatusPendingSync
		if ev.Deleted {
			status = domain.ContactStatusDeletedInAmo
		}

		contactID, err := u.repo.UpsertContactFromWebhook(accountID, ev.AmoID, ev.Name, email, status)
		if err != nil {
			return 0, err
		}

		if err = u.repo.AddSyncHistory(contactID, "webhook_received", "webhook received"); err != nil {
			log.Printf("contact_id=%d err=%v", contactID, err)
		}
		if err = u.repo.TrimSyncHistory(contactID, 10); err != nil {
			log.Printf("contact_id=%d err=%v", contactID, err)
		}

		if ev.Deleted {
			deleteIDs = append(deleteIDs, ev.AmoID)
			continue
		}
		upsertIDs = append(upsertIDs, ev.AmoID)
	}

	var lastJobID uint64

	if len(upsertIDs) > 0 {
		id, err := u.producer.PushWebhookUpsertJob(ctx, accountID, upsertIDs)
		if err != nil {
			return 0, err
		}
		lastJobID = id

		for _, amoID := range upsertIDs {
			contactID, ok, err := u.repo.FindContactIDByAmoID(accountID, amoID)
			if err != nil || !ok {
				continue
			}
			if err = u.repo.AddSyncHistory(contactID, "queued", "enqueued to beanstalk (upsert)"); err != nil {
				log.Printf("contact_id=%d err=%v", contactID, err)
			}
			if err = u.repo.TrimSyncHistory(contactID, 10); err != nil {
				log.Printf("contact_id=%d err=%v", contactID, err)
			}
		}
	}

	if len(deleteIDs) > 0 {
		id, err := u.producer.PushWebhookDeleteJob(ctx, accountID, deleteIDs)
		if err != nil {
			return 0, err
		}
		lastJobID = id

		for _, amoID := range deleteIDs {
			contactID, ok, err := u.repo.FindContactIDByAmoID(accountID, amoID)
			if err != nil || !ok {
				continue
			}
			if err = u.repo.AddSyncHistory(contactID, "queued", "enqueued to beanstalk (delete)"); err != nil {
				log.Printf("contact_id=%d err=%v", contactID, err)
			}
			if err = u.repo.TrimSyncHistory(contactID, 10); err != nil {
				log.Printf("contact_id=%d err=%v", contactID, err)
			}
		}
	}

	return lastJobID, nil
}

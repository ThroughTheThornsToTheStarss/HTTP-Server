package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/queue"
	qbs "git.amocrm.ru/ilnasertdinov/http-server-go/internal/queue/beanstalk"
	mysqlrepo "git.amocrm.ru/ilnasertdinov/http-server-go/internal/repo/mysql"
	mysqlcfg "git.amocrm.ru/ilnasertdinov/http-server-go/pkg/mysql"

	"github.com/beanstalkd/go-beanstalk"
)

const (
	workerKindsSetCap = 8
	workerKindsSep    = ","
)

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	beanstalkAddr := getenv("BEANSTALK_ADDR", "beanstalkd:11300")

	consumer, err := qbs.NewConsumer(beanstalkAddr)
	if err != nil {
		log.Fatalf("beanstalk consumer init: %v", err)
	}
	defer func() { _ = consumer.Close() }()

	db, err := mysqlcfg.NewGormFromEnv()
	if err != nil {
		log.Fatalf("mysql connect error: %v", err)
	}

	repo := mysqlrepo.NewGormRepository(db)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	allowedKinds := allowedKindsFromEnv()
	if allowedKinds == nil {
		log.Printf("worker kinds: all")
	} else {
		log.Printf("worker kinds: %v", allowedKinds)
	}

	log.Printf("worker started; beanstalk=%s", beanstalkAddr)

	for {
		select {
		case <-ctx.Done():
			log.Printf("worker stopped")
			return
		default:
		}

		jobID, job, err := consumer.Reserve(ctx, 5*time.Second)
		if err != nil {
			if errors.Is(err, beanstalk.ErrTimeout) {
				continue
			}
			log.Printf("reserve error: %v", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if !kindAllowed(allowedKinds, job.Kind) {
			_ = consumer.Release(jobID, 2*time.Second)
			continue
		}

		switch job.Kind {
		case queue.JobKindInitialSync, queue.JobKindWebhookUpsert, queue.JobKindWebhookDelete:
		default:
			log.Printf("unknown job kind: id=%d kind=%q -> bury", jobID, job.Kind)
			_ = consumer.Bury(jobID)
			continue
		}

		if job.Kind == queue.JobKindWebhookUpsert || job.Kind == queue.JobKindWebhookDelete {
			for _, amoID := range job.ContactIDs {
				cid, ok, err := repo.FindContactIDByAmoID(job.AccountID, amoID)
				if err != nil || !ok {
					continue
				}
				_ = repo.AddSyncHistory(cid, "picked", "worker picked job")
				_ = repo.TrimSyncHistory(cid, 10)
			}
		}

		log.Printf("job accepted: id=%d kind=%s account_id=%d contact_ids=%v", jobID, job.Kind, job.AccountID, job.ContactIDs)

		if err := processJob(ctx, repo, job); err != nil {
			if errors.Is(err, errNonRetryable) {
				log.Printf("process error: id=%d kind=%s err=%v -> bury", jobID, job.Kind, err)
				_ = consumer.Bury(jobID)
			} else {
				log.Printf("process error: id=%d kind=%s err=%v -> release", jobID, job.Kind, err)
				_ = consumer.Release(jobID, 10*time.Second)
			}
			continue
		}

		if err := consumer.Delete(jobID); err != nil {
			log.Printf("delete job error: id=%d err=%v", jobID, err)
		}
	}
}

var errNonRetryable = errors.New("non-retryable job")

type contactRepo interface {
	FindContactIDByAmoID(accountID uint64, amoID int64) (uint, bool, error)
	AddSyncHistory(contactID uint, status string, message string) error
	TrimSyncHistory(contactID uint, keepLast int) error
}

func processJob(ctx context.Context, repo contactRepo, job queue.Job) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	switch job.Kind {
	case queue.JobKindInitialSync:
		return nil

	case queue.JobKindWebhookUpsert, queue.JobKindWebhookDelete:
		if len(job.ContactIDs) == 0 {
			return fmt.Errorf("%w: empty contact_ids", errNonRetryable)
		}

		for _, amoID := range job.ContactIDs {
			cid, ok, err := repo.FindContactIDByAmoID(job.AccountID, amoID)
			if err != nil {
				return err
			}
			if !ok {
				continue
			}
			if err := repo.AddSyncHistory(cid, "done", "worker processed job"); err != nil {
				return err
			}
			if err := repo.TrimSyncHistory(cid, 10); err != nil {
				return err
			}
		}
		return nil

	default:
		return fmt.Errorf("%w: unknown kind %q", errNonRetryable, job.Kind)
	}
}
func allowedKindsFromEnv() map[string]struct{} {
	raw := strings.TrimSpace(os.Getenv("WORKER_KINDS"))
	if raw == "" {
		return nil
	}

	set := make(map[string]struct{}, workerKindsSetCap)
	for _, p := range strings.Split(raw, workerKindsSep) {
		k := strings.TrimSpace(p)
		if k != "" {
			set[k] = struct{}{}
		}
	}
	return set
}

func kindAllowed(set map[string]struct{}, kind string) bool {
	if set == nil {
		return true
	}
	_, ok := set[kind]
	return ok
}

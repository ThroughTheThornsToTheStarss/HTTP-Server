package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/queue"
	qbs "git.amocrm.ru/ilnasertdinov/http-server-go/internal/queue/beanstalk"
	mysqlrepo "git.amocrm.ru/ilnasertdinov/http-server-go/internal/repo/mysql"
	mysqlcfg "git.amocrm.ru/ilnasertdinov/http-server-go/pkg/mysql"

	"github.com/beanstalkd/go-beanstalk"
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

		if err := consumer.Delete(jobID); err != nil {
			log.Printf("delete job error: id=%d err=%v", jobID, err)
		}

		if job.Kind == queue.JobKindWebhookUpsert || job.Kind == queue.JobKindWebhookDelete {
			for _, amoID := range job.ContactIDs {
				cid, ok, err := repo.FindContactIDByAmoID(job.AccountID, amoID)
				if err != nil || !ok {
					continue
				}
				_ = repo.AddSyncHistory(cid, "done", "worker acked job")
				_ = repo.TrimSyncHistory(cid, 10)
			}
		}
	}
}

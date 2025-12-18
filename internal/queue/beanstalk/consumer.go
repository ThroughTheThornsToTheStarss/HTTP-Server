package beanstalk

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/queue"
	"github.com/beanstalkd/go-beanstalk"
)

type Consumer struct {
	conn    *beanstalk.Conn
	tubeSet *beanstalk.TubeSet
}

func NewConsumer(addr string) (*Consumer, error) {
	if addr == "" {
		return nil, errors.New("beanstalk addr is empty")
	}

	conn, err := beanstalk.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		conn:    conn,
		tubeSet: beanstalk.NewTubeSet(conn, tubeName),
	}, nil
}

func (c *Consumer) Reserve(ctx context.Context, timeout time.Duration) (uint64, queue.Job, error) {
	if err := ctx.Err(); err != nil {
		return 0, queue.Job{}, err
	}

	id, body, err := c.tubeSet.Reserve(timeout)
	if err != nil {
		return 0, queue.Job{}, err
	}

	var job queue.Job
	if err := json.Unmarshal(body, &job); err != nil {
		_ = c.conn.Bury(id, defaultPriority)
		return 0, queue.Job{}, err
	}

	return id, job, nil
}

func (c *Consumer) Delete(id uint64) error { return c.conn.Delete(id) }

func (c *Consumer) Release(id uint64, delay time.Duration) error {
	return c.conn.Release(id, defaultPriority, delay)
}

func (c *Consumer) Bury(id uint64) error { return c.conn.Bury(id, defaultPriority) }

func (c *Consumer) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

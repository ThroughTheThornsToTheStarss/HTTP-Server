package beanstalk

import (
	"errors"
	"time"

	"git.amocrm.ru/ilnasertdinov/http-server-go/internal/queue"
	"github.com/beanstalkd/go-beanstalk"
)

const (
	tubeName        = "sync_contacts"
	defaultPriority = 1024
	defaultDelay    = 0 * time.Second
	defaultTTR      = 120 * time.Second
)

type Producer struct {
	conn *beanstalk.Conn
	tube *beanstalk.Tube
}

var _ queue.Producer = (*Producer)(nil)

func New(addr string) (*Producer, error) {
	if addr == "" {
		return nil, errors.New("beanstalk addr is empty")
	}

	conn, err := beanstalk.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Producer{
		conn: conn,
		tube: &beanstalk.Tube{Conn: conn, Name: tubeName},
	}, nil
}

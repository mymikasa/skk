package events

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/mymikasa/skk/interactive/repository"
	"github.com/mymikasa/skk/pkg/logger"
	"github.com/mymikasa/skk/pkg/saramax"
	"time"
)

const TopicReadEvents = "prompt_read"

type InteractiveReadEventConsumer struct {
	repo   repository.InteractiveRepository
	client sarama.Client
	l      logger.LoggerV1
}

func NewInteractiveReadEventConsumer(repo repository.InteractiveRepository, client sarama.Client, l logger.LoggerV1) *InteractiveReadEventConsumer {
	return &InteractiveReadEventConsumer{repo: repo, client: client, l: l}
}

func (i *InteractiveReadEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", i.client)

	if err != nil {
		return err
	}

	go func() {
		er := cg.Consume(context.Background(), []string{TopicReadEvents},
			saramax.NewHandler[ReadEvent](i.l, i.Consume))

		if er != nil {
			i.l.Error("退出消费", logger.Error(er))
		}
	}()
	return err
}

type ReadEvent struct {
	Pid int64
	Uid int64
}

func (i *InteractiveReadEventConsumer) Consume(msg *sarama.ConsumerMessage, event ReadEvent) error {
	ctx, cancle := context.WithTimeout(context.Background(), time.Second)
	defer cancle()

	return i.repo.IncrReadCnt(ctx, "prompt", event.Pid)
}

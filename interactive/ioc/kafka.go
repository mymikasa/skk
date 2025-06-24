package ioc

import (
	"github.com/IBM/sarama"
	events2 "github.com/mymikasa/skk/interactive/events"
	"github.com/mymikasa/skk/pkg/events"
	"github.com/spf13/viper"
)

func InitSaramaClient() sarama.Client {
	type Confif struct {
		Addr []string `yaml:"addr"`
	}

	var cfg Confif

	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}

	scfg := sarama.NewConfig()
	scfg.Producer.Return.Successes = true
	client, err := sarama.NewClient(cfg.Addr, scfg)

	if err != nil {
		panic(err)
	}
	return client
}

func InitConsumers(c1 *events2.InteractiveReadEventConsumer) []events.Consumer {
	return []events.Consumer{c1}
}

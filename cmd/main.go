package main

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/go1com/es-writer"
)

func main() {
	ctx := context.Background()
	f := es_writer.NewFlags()

	// Credentials can be leaked with debug enabled.
	if *f.Debug {
		logrus.Infoln("======= ElasticSearch-Writer =======")
		logrus.Infof("RabbitMQ URL: %s", *f.Url)
		logrus.Infof("RabbitMQ kind: %s", *f.Kind)
		logrus.Infof("RabbitMQ exchange: %s", *f.Exchange)
		logrus.Infof("RabbitMQ routing key: %s", *f.RoutingKey)
		logrus.Infof("RabbitMQ prefetch count: %d", *f.PrefetchCount)
		logrus.Infof("RabbitMQ prefetch size: %d", *f.PrefetchSize)
		logrus.Infof("RabbitMQ queue name: %s", *f.QueueName)
		logrus.Infof("RabbitMQ consumer name: %s", *f.ConsumerName)
		logrus.Infof("ElasticSearch URL: %s", *f.EsUrl)
		logrus.Infof("Tick interval: %s", *f.TickInterval)
		logrus.Infoln("====================================")

		logrus.SetLevel(logrus.DebugLevel)
	}

	dog, err, stop := f.Dog()
	if err != nil {
		logrus.
			WithError(err).
			Panicln("failed to get the dog")
	}

	defer func() { stop <- true }()

	go es_writer.StartPrometheusServer(*f.AdminPort)
	dog.Start(ctx, f)

	os.Exit(1)
}

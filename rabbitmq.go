package a_wip

import (
	"github.com/streadway/amqp"
	"github.com/Sirupsen/logrus"
)

func Connection(queueUrl string) (*amqp.Connection, error) {
	con, err := amqp.Dial(queueUrl)

	if nil != err {
		return nil, err
	}

	go func() {
		conCloseChan := con.NotifyClose(make(chan *amqp.Error))

		select
		{
		case err := <-conCloseChan:
			logrus.Errorln(err.Error())

			panic(err.Error())
		}
	}()

	return con, nil
}

func Channel(con *amqp.Connection, kind string, exchangeName string, prefetchCount int, prefetchSize int) *amqp.Channel {
	ch, err := con.Channel()

	if nil != err {
		panic(err.Error())
	}

	if "topic" != kind && "direct" != kind {
		panic("Unsupported channel Kind: " + kind)
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		kind,
		false,
		false,
		false,
		false,
		nil,
	)

	if nil != err {
		panic(err.Error())
	}

	err = ch.Qos(prefetchCount, prefetchSize, false)
	if nil != err {
		panic(err.Error())
	}

	return ch
}

func Messages(ch *amqp.Channel, queueName string, exchange string, routingKey string, consumerName string) <-chan amqp.Delivery {
	queue, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)

	if nil != err {
		panic(err.Error())
	}

	ch.QueueBind(queue.Name, routingKey, exchange, true, nil)

	messages, err := ch.Consume(
		queue.Name,
		consumerName,
		false,
		false,
		false,
		true,
		nil,
	)

	if nil != err {
		panic(err.Error())
	}

	return messages
}

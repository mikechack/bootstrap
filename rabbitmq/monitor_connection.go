package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
)

func monitorConnection(connection *amqp.Connection) {
	go doMonitor(connection)
}

func monitorChannel(channel *amqp.Channel) {
	go doMonitorChannel(channel)
}

func doMonitor(connection *amqp.Connection) {
	c := make(chan *amqp.Error)
	connection.NotifyClose(c)
	select {
	case err := <-c:
		log.Printf("RabbitMQ Connection Error - %s\n\n", err)
	}

}

func doMonitorChannel(channel *amqp.Channel) {
	c := make(chan *amqp.Error)
	channel.NotifyClose(c)
	select {
	case err := <-c:
		log.Printf("RabbitMQ Channel Error - %s\n\n", err)
	}

}

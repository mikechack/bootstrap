package rabbitmq

import (
	"bytes"
	"crypto/tls"
	"github.com/mikechack/bootstrap/certs"
	"github.com/mikechack/bootstrap/models"
	"github.com/mikechack/bootstrap/utility"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type rmqInfo struct {
	connectorUuid        string
	exchangeRegistration string
	exchangeLocal        string
	topicRegJoin         string
	localRegQueue        string
	localRegTopic        string
	channel              *amqp.Channel
	connection           *amqp.Connection
}

func (info *rmqInfo) init(connectorUuid, exchgReg, organization, topicRegJ string) {
	info.connectorUuid = connectorUuid
	info.exchangeRegistration = exchgReg
	info.exchangeLocal = "mediafusion:" + organization
	info.topicRegJoin = topicRegJ
	info.localRegQueue = "register." + connectorUuid
	info.localRegTopic = "register." + connectorUuid
}

var rmqContext rmqInfo = rmqInfo{}

var cert tls.Certificate
var connection *amqp.Connection
var channel *amqp.Channel

func createUrl(name, password string) string {
	var buffer bytes.Buffer

	buffer.WriteString("amqps://")
	buffer.WriteString(name)
	buffer.WriteString(":")
	buffer.WriteString(password)
	buffer.WriteString("@sj21lab-rabbitmq-1.cisco.com:5671/")

	return buffer.String()
}

func listenAmqp() error {
	var err error

	log.Printf("Device exchange name %s\n", rmqContext.exchangeLocal)

	durable := false
	autoDelete, noWait := false, false
	internal := false
	channel, _ := rmqContext.connection.Channel()
	tries := 0
	timer := time.NewTicker(1 * time.Second)
	for _ = range timer.C {
		err = channel.ExchangeDeclarePassive(rmqContext.exchangeLocal, "topic", durable, autoDelete, internal, noWait, nil)
		if err != nil {
			channel, _ = rmqContext.connection.Channel()
			log.Printf("Could not declare personal queue for device %v\n", err)
			tries++
			if tries > 5 {
				break
			}
			continue
		}
		log.Printf("Personal queue is a go")
		err = nil
		break
	}
	timer.Stop()
	channel.Close()

	return err
}

func ConnectAmqp() error {

	var err error
	//cert, err = tls.LoadX509KeyPair("certs/client0-ca.crt", "certs/client0.key")
	clientCert, _ := certs.Asset("certs/client0.crt")
	clientKey, _ := certs.Asset("certs/client0.key")
	cert, err = tls.X509KeyPair(clientCert, clientKey)
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	session := models.GetSession()
	url := createUrl(session.Ma_name, session.Access_token)
	log.Printf("AMQP url = %s\n", url)
	connection, err = amqp.DialTLS(url, &tls.Config{Certificates: []tls.Certificate{cert}, PreferServerCipherSuites: true, InsecureSkipVerify: true})
	rmqContext.connection = connection
	monitorConnection(connection)
	models.ReturnSession(session)

	connectorId := utility.GetRandomString(16)
	rmqContext.init(connectorId, "mediafusion:device.registration", session.Ma_organization, "register.join")

	durable := false
	autoDelete, noWait := false, false
	internal := false
	channel, err := connection.Channel()
	if err != nil {
		log.Printf("Creating channel was an issue: %s", err)
		return err
	}
	monitorChannel(channel)
	rmqContext.channel = channel
	err = channel.ExchangeDeclarePassive(rmqContext.exchangeRegistration, "topic", durable, autoDelete, internal, noWait, nil)
	if err != nil {
		log.Printf("Exchange does not exist: %s", err)
		return err
	}

	//sendRegJoin("12341234")

	//sendRegStatus("reguuid", rmqContext.topicRegJoin)

	if err = listenAmqp(); err != nil {
		return err
	}

	go registrationListener(channel)

	return err
}

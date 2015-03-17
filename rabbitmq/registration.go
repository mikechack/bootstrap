package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"sqbu.com/MediaFusion/bootstrap/models"
)

type regReport struct {
	Filter []regDevice `json:"filter,omitempty"`
}

type regDevice struct {
	Reguuid  string `json:"reguuid,omitempty"`
	Deviceid string `json:"deviceid,omitempty"`
}

type regStatusHolder struct {
	Status regStatus `json:"status,omitempty"`
}

type regStatus struct {
	Reguuid string `json:"reguuid,omitempty"`
	State   string `json:"state,omitempty"`
}

func SendRegJoin(serialNumber string) {
	session := models.GetSession()
	models.ReturnSession(session)
	dev, _ := models.GetDeviceBySerialNumber(serialNumber)
	dev.ConnectorUuid = rmqContext.connectorUuid
	dev.OrganizationId = session.Ma_organization
	log.Printf("SendRegJoin - Device SerialNumber is %s", dev.SerialNumber)

	buf, _ := json.Marshal(dev)

	msg := amqp.Publishing{
		DeliveryMode: 1,
		//Timestamp:    t,
		ContentType: "application/registration.v1+json",
		Body:        buf,
		Type:        "REG_JOIN",
	}
	mandatory, immediate := false, false
	rmqContext.channel.Publish(rmqContext.exchangeRegistration, rmqContext.topicRegJoin, mandatory, immediate, msg)
}

func sendRegStatus(reguuid string, destination string) {
	dev, _ := models.GetDeviceByReguuid(reguuid)
	log.Printf("SendRegStatus - Device SerialNumber is %s", dev.SerialNumber)

	status := regStatusHolder{}
	status.Status = regStatus{Reguuid: reguuid, State: "register"}
	buf, _ := json.Marshal(status)

	msg := amqp.Publishing{
		DeliveryMode: 1,
		//Timestamp:    t,
		ContentType: "application/registration.v1+json",
		Body:        buf,
		Type:        "REG_STATUS",
	}
	mandatory, immediate := false, false
	rmqContext.channel.Publish(rmqContext.exchangeRegistration, destination, mandatory, immediate, msg)
}

func unmarshallRegReport(message []byte) (regReport, error) {
	var request regReport
	var err error
	err = json.Unmarshal(message[0:], &request)
	if err != nil {
		log.Printf("Unmarshall Reg Report Error %s\n", err.Error())
		return request, err
	} else {
		log.Printf("Unmarshall Reg Report Good - reguuid[0] = %s\n", request.Filter[0].Reguuid)
		models.MapDeviceByReguuid(request.Filter[0].Reguuid, request.Filter[0].Deviceid)
		sendRegStatus(request.Filter[0].Reguuid, rmqContext.topicRegJoin)
	}
	return request, err
}

func registrationListener(channel *amqp.Channel) {

	durable, exclusive := false, false
	autoDelete, noWait := true, true
	log.Printf("Registration listener -  %s   -  %s\n", rmqContext.localRegQueue, rmqContext.exchangeLocal)
	//queueName := "register." + connectionId
	q, _ := channel.QueueDeclare(rmqContext.localRegQueue, durable, autoDelete, exclusive, noWait, nil)
	channel.QueueBind(q.Name, rmqContext.localRegTopic, rmqContext.exchangeLocal, false, nil)
	autoAck, exclusive, noLocal, noWait := false, false, false, false
	messages, _ := channel.Consume(q.Name, "", autoAck, exclusive, noLocal, noWait, nil)
	multiAck := false
	for msg := range messages {
		log.Printf("Register Type Message\n\tReceived Body:   %s\n\tContent Type:   %s\n", string(msg.Body), msg.ContentType)
		unmarshallRegReport(msg.Body)
		msg.Ack(multiAck)
	}
	log.Printf("Registration listener exited\n")

}

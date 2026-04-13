package messaging

import (
	"encoding/json"
	"log"
	"ride-sharing/shared/contracts"
)

type QueueConsumer struct {
	rb        *RabbitMQ
	connMng   *ConnectionManager
	queueName string
	transform func(routingKey string, data []byte) (any, error)
}

func NewQueueConsumer(rb *RabbitMQ, connMngr *ConnectionManager, queueName string) *QueueConsumer {
	return &QueueConsumer{
		rb:        rb,
		connMng:   connMngr,
		queueName: queueName,
	}
}

func NewQueueConsumerWithTransform(rb *RabbitMQ, connMngr *ConnectionManager, queueName string, transform func(routingKey string, data []byte) (any, error)) *QueueConsumer {
	return &QueueConsumer{
		rb:        rb,
		connMng:   connMngr,
		queueName: queueName,
		transform: transform,
	}
}

func (qc *QueueConsumer) Start() error {
	msgs, err := qc.rb.Channel.Consume(
		qc.queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var msgBody contracts.AmqpMessage
			if err := json.Unmarshal(msg.Body, &msgBody); err != nil {
				log.Println("Failed to unmarshal message:", err)
			}

			userID := msgBody.OwnerID

			var payload any
			if msgBody.Data != nil {
				if qc.transform != nil {
					payload, err = qc.transform(msg.RoutingKey, msgBody.Data)
					if err != nil {
						log.Println("Failed to transform payload:", err)
						continue
					}
				} else {
					if err := json.Unmarshal(msgBody.Data, &payload); err != nil {
						log.Println("Failed to unmarshal payload:", err)
						continue
					}
				}
			}

			clientMsg := contracts.WSMessage{
				Type: msg.RoutingKey,
				Data: payload,
			}

			if err := qc.connMng.SendMessage(userID, clientMsg); err != nil {
				log.Printf("Failed to send message to user %s: %v", userID, err)
			}
		}
	}()

	return nil

}

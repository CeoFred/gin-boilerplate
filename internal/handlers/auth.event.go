package handlers

import (
	"encoding/json"
	"fmt"
	// "log"

	bootsreap "github.com/CeoFred/gin-boilerplate/internal/bootstrap"
	"github.com/CeoFred/gin-boilerplate/internal/models"

	"github.com/IBM/sarama"
)

type EventHandler struct {
	Deps *bootsreap.AppDependencies
}

func (h *EventHandler) ProcessSignup(msg *sarama.ConsumerMessage) error {
	// log.Printf("Processing message: topic=%s, partition=%d, offset=%d, key=%s, value=%s",
	// 	msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
	user := &models.User{}
	if err := json.Unmarshal(msg.Value, user); err != nil {
		return err
	}

	fmt.Println("new user registration event received", user)

	// maybe send email notification
	return nil
}

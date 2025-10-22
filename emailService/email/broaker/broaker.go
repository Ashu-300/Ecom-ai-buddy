package broaker

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"supernova/emailService/email/controller"

	"github.com/streadway/amqp"
)

type JsonUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

var (
	conn       *amqp.Connection
	ch         *amqp.Channel
	notifyClose chan *amqp.Error
	mutex      sync.Mutex
)

// Connect sets up a RabbitMQ consumer that auto-reconnects and processes durable messages
func Connect() {
	for {
		err := connectAndConsume()
		if err != nil {
			log.Println("⚠️ RabbitMQ connection lost or failed:", err)
		}
		log.Println("🔁 Reconnecting to RabbitMQ in 5s...")
		time.Sleep(5 * time.Second)
	}
}

func connectAndConsume() error {
	mutex.Lock()
	defer mutex.Unlock()

	amqpServerURL := os.Getenv("AMQP_SERVER_URL")
	if amqpServerURL == "" {
		log.Fatal("❌ Missing AMQP_SERVER_URL environment variable")
	}

	var err error
	conn, err = amqp.Dial(amqpServerURL)
	if err != nil {
		log.Println("❌ Failed to connect to RabbitMQ:", err)
		return err
	}

	ch, err = conn.Channel()
	if err != nil {
		log.Println("❌ Failed to open channel:", err)
		return err
	}
	queues := []string{"AuthService", "PaymentService",}

	for _, q := range queues {
		_, err := ch.QueueDeclare(
			q,    // queue name
			true, // durable
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Printf("❌ Failed to declare queue %s: %v", q, err)
			return err
		}
	}

	// Set prefetch to 1 to process one message at a time
	if err := ch.Qos(1, 0, false); err != nil {
		log.Println("⚠️ Failed to set QoS:", err)
	}


	log.Println("✅ Connected to RabbitMQ (Consumer)")
	log.Println("📩 Waiting for messages...")

	// Listen for unexpected close events
	notifyClose = make(chan *amqp.Error)
	ch.NotifyClose(notifyClose)

	// Start message processing loop

	for _, q := range queues {
		msgs, _ := ch.Consume(q, "", false, false, false, false, nil)
		go func(queue string, msgs <-chan amqp.Delivery) {
			for msg := range msgs {
				handleMessage(queue, msg)
			}
		}(q, msgs)
	}
	log.Println("✅ Connected to RabbitMQ and consuming multiple queues")
	return nil
}

func handleMessage(queue string, msg amqp.Delivery) {
    switch queue {
    case "AuthService":
        var user JsonUser
        json.Unmarshal(msg.Body, &user)
        controller.SendEmail(user.Email, user.Name)
    // case "PaymentService":
    //     var data ResetData
    //     json.Unmarshal(msg.Body, &data)
    //     controller.SendResetEmail(data.Email, data.Token)
    // case "Promotions":
    //     var promo PromoData
    //     json.Unmarshal(msg.Body, &promo)
    //     controller.SendPromoEmail(promo.Email, promo.Offer)
    }
    msg.Ack(false)
}


func GetChannel() *amqp.Channel {
	mutex.Lock()
	defer mutex.Unlock()
	return ch
}

func GetConnection() *amqp.Connection {
	mutex.Lock()
	defer mutex.Unlock()
	return conn
}

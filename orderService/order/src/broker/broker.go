package broker

import (
	"log"
	"os"
	"sync"
	"time"

	

	"github.com/streadway/amqp"
)

var (
	conn        *amqp.Connection
	channel     *amqp.Channel
	notifyClose chan *amqp.Error
	mutex       sync.Mutex
	amqpURL     string
	retryBackoff = 5 * time.Second
)


// Connect initializes RabbitMQ connection and channel (idempotent)
func Connect() {
	mutex.Lock()
	defer mutex.Unlock()

	if amqpURL == "" {
		amqpURL = os.Getenv("AMQP_SERVER_URL")
		if amqpURL == "" {
			log.Fatal("❌ sellerDashboard AMQP_SERVER_URL not set")
		}
	}

	for {
		var err error
		log.Println("🔁 sellerDashboard service Connecting to RabbitMQ...")
		conn, err = amqp.Dial(amqpURL)
		if err != nil {
			log.Println("⚠️ sellerDashboard Failed to connect:", err)
			time.Sleep(retryBackoff)
			continue
		}

		channel, err = conn.Channel()
		if err != nil {
			log.Println("⚠️ sellerDashboard Failed to open channel:", err)
			_ = conn.Close()
			time.Sleep(retryBackoff)
			continue
		}

		notifyClose = make(chan *amqp.Error)
		channel.NotifyClose(notifyClose)

		// Launch reconnect handler in background
		go handleReconnect(notifyClose)

		log.Println("✅ sellerDashboard service  Connected to RabbitMQ")
		return
	}
}

func handleReconnect(nc chan *amqp.Error) {
	err := <-nc
	if err != nil {
		log.Printf("🚨 sellerDashboard RabbitMQ closed: %v. Reconnecting...", err)
	} else {
		log.Println("ℹ️ sellerDashboard RabbitMQ NotifyClose returned nil. Reconnecting...")
	}

	mutex.Lock()
	if channel != nil {
		_ = channel.Close()
	}
	if conn != nil {
		_ = conn.Close()
	}
	channel, conn = nil, nil
	mutex.Unlock()

	// reconnect in background
	for {
		Connect()
		mutex.Lock()
		ok := conn != nil && channel != nil
		mutex.Unlock()
		if ok {
			log.Println("✅ sellerDashboard Reconnected to RabbitMQ (background)")
			return
		}
		time.Sleep(retryBackoff)
	}
}

// PublishJSON sends a persistent JSON message to a queue
func PublishJSON(queueName string, body []byte) error {
	if conn == nil || channel == nil {
		Connect()
	}

	mutex.Lock()
	ch := channel
	mutex.Unlock()
	if ch == nil {
		return amqp.ErrClosed
	}
	_, err := ch.QueueDeclare(
		queueName,  // queue name
		true,       // durable
		false,      // auto-delete
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"", queueName, false, false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		log.Println("❌ sellerDashboard Publish failed, reconnecting...")
		Connect()
		mutex.Lock()
		ch = channel
		mutex.Unlock()
		if ch == nil {
			return amqp.ErrClosed
		}
		return ch.Publish(
			"", queueName, false, false,
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         body,
				DeliveryMode: amqp.Persistent,
				Timestamp:    time.Now(),
			},
		)
	}

	log.Printf("📤 sellerDashboard Sent message to %s", queueName)
	return nil
}

// ConsumeQueues sets up consumers for multiple queues
// func ConsumeQueues() {
// 	if conn == nil || channel == nil {
// 		Connect()
// 	}

// 	for _, q := range queues {
// 		_, err := channel.QueueDeclare(q, true, false, false, false, nil)
// 		if err != nil {
// 			log.Fatalf("❌ sellerDashboard Failed to declare queue %s: %v", q, err)
// 		}

// 		msgs, err := channel.Consume(q, "", false, false, false, false, nil)
// 		if err != nil {
// 			log.Fatalf("❌ sellerDashboard Failed to consume queue %s: %v", q, err)
// 		}

// 		go func(queue string, msgs <-chan amqp.Delivery) {
// 			for msg := range msgs {
// 				handleMessage(queue, msg)
// 			}
// 		}(q, msgs)
// 	}
// 	log.Println("✅ sellerdashboard Consumers started for queues:", queues)
// }

// func handleMessage(queue string, msg amqp.Delivery) {
// 	switch queue {
// 	case "AuthService":
// 		var user dto.JsonUser
// 		_ = json.Unmarshal(msg.Body, &user)
// 		// controller.AuthEmail(user.Email, user.Name)
// 	case "PaymentService":
// 		var data dto.PaymentData
// 		_ = json.Unmarshal(msg.Body, &data)
// 		// controller.PaymentSuccessEmail(data)
// 	case "AuthServiceDashboard":
// 		var user models.User
// 		_ = json.Unmarshal(msg.Body, &user)
// 		 controller.CreateUser(user)
// 	case "ProductDashboard":
// 		var product models.Product
// 		_ = json.Unmarshal(msg.Body , &product)
// 		controller.CreateProduct(product)

// 	}
// 	msg.Ack(false)
// }

// GetChannel returns current channel
func GetChannel() *amqp.Channel {
	mutex.Lock()
	defer mutex.Unlock()
	return channel
}

// GetConnection returns current connection
func GetConnection() *amqp.Connection {
	mutex.Lock()
	defer mutex.Unlock()
	return conn
}

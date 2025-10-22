package broker

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

var (
	conn         *amqp.Connection
	channel      *amqp.Channel
	notifyClose  chan *amqp.Error
	mutex        sync.Mutex
	amqpURL      string
	retryBackoff = 5 * time.Second
)

// ConnectBroker establishes connection and channel to RabbitMQ and declares the durable queue.
// It is safe to call multiple times (uses a mutex).
func ConnectBroker() {
	mutex.Lock()
	defer mutex.Unlock()

	// read AMQP url once
	if amqpURL == "" {
		amqpURL = os.Getenv("AMQP_SERVER_URL")
		if amqpURL == "" {
			log.Fatal("AMQP_SERVER_URL is not set")
		}
	}

	// If already connected and channel seems fine, just return
	if conn != nil && channel != nil {
		return
	}

	for {
		var err error
		log.Println("🔁 Attempting to connect to RabbitMQ...")
		conn, err = amqp.Dial(amqpURL)
		if err != nil {
			log.Println("⚠️  Failed to connect to RabbitMQ:", err)
			time.Sleep(retryBackoff)
			continue
		}

		channel, err = conn.Channel()
		if err != nil {
			log.Println("⚠️  Failed to open channel:", err)
			_ = conn.Close()
			conn = nil
			time.Sleep(retryBackoff)
			continue
		}

		// Declare a durable queue (idempotent; safe to declare from both producer and consumer)
		_, err = channel.QueueDeclare(
			"PaymentService", // queue name
			true,          // durable
			false,         // auto-delete
			false,         // exclusive
			false,         // no-wait
			nil,           // args
		)
		if err != nil {
			log.Println("⚠️  Failed to declare queue:", err)
			_ = channel.Close()
			_ = conn.Close()
			channel = nil
			conn = nil
			time.Sleep(retryBackoff)
			continue
		}

		// Enable publisher confirms - best-effort
		if err := channel.Confirm(false); err != nil {
			log.Println("⚠️  Could not put channel into confirm mode:", err)
			// not fatal; we continue
		}

		// Setup NotifyClose to detect unexpected channel/connection closures
		notifyClose = make(chan *amqp.Error)
		channel.NotifyClose(notifyClose)

		// Launch a goroutine to handle closures and attempt reconnect
		go func(nc chan *amqp.Error) {
			err := <-nc
			if err != nil {
				log.Printf("🚨 RabbitMQ channel closed: %v. Will attempt reconnect...\n", err)
			} else {
				log.Println("ℹ️ RabbitMQ NotifyClose returned nil error (channel closed). Reconnecting...")
			}
			// Clean up existing references; next publish will call ConnectBroker again
			mutex.Lock()
			if channel != nil {
				_ = channel.Close()
			}
			if conn != nil {
				_ = conn.Close()
			}
			channel = nil
			conn = nil
			mutex.Unlock()

			// Try to reconnect in background (non-blocking)
			for {
				ConnectBroker()
				// if connection succeeded, break
				mutex.Lock()
				ok := (conn != nil && channel != nil)
				mutex.Unlock()
				if ok {
					log.Println("✅ Reconnected to RabbitMQ (background)")
					return
				}
				time.Sleep(retryBackoff)
			}
		}(notifyClose)

		log.Println("✅ Successfully connected to RabbitMQ (Producer)")
		return
	}
}

// PublishJSON sends a persistent JSON message to the specified queue.
// On failure it attempts a reconnect and retries once.
func PublishJSON(queueName string, body []byte) error {
	// Ensure we have a connection & channel
	if conn == nil || channel == nil {
		ConnectBroker()
	}

	// mutex to avoid racing ConnectBroker from NotifyClose goroutine
	mutex.Lock()
	ch := channel
	mutex.Unlock()

	if ch == nil {
		// if still nil after ConnectBroker, return an error
		return amqp.ErrClosed
	}

	// publish attempt
	err := ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // persistent messages
			Timestamp:    time.Now(),
		},
	)
	if err == nil {
		// attempt to wait for confirm (if channel was put into confirm mode)
		ackCh := ch.NotifyPublish(make(chan amqp.Confirmation, 1))
		select {
		case conf := <-ackCh:
			if conf.Ack {
				log.Printf("📤 Sent message to %s: %s", queueName, string(body))
				return nil
			}
			log.Println("❌ Message not acknowledged by broker")
			// fallthrough to retry logic
		case <-time.After(5 * time.Second):
			// confirmation timeout - treat as possible failure but do not block forever
			log.Println("⚠️ Timeout waiting for broker confirmation (continuing)")
			// We still treat publish as success because Publish didn't return an error.
			log.Printf("📤 Sent (no confirm) message to %s: %s", queueName, string(body))
			return nil
		}
	} else {
		log.Println("❌ Failed to publish (first attempt):", err)
	}

	// If we reached here -> first publish failed or not acked. Try reconnect + retry once.
	log.Println("🔁 Attempting reconnect and single retry...")
	ConnectBroker()

	// get fresh channel
	mutex.Lock()
	ch = channel
	mutex.Unlock()
	if ch == nil {
		return amqp.ErrClosed
	}

	err = ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		log.Println("🚨 Retry publish failed:", err)
		return err
	}

	// attempt confirm again (best-effort)
	ackCh := ch.NotifyPublish(make(chan amqp.Confirmation, 1))
	select {
	case conf := <-ackCh:
		if conf.Ack {
			log.Printf("📤 Sent message to %s (after retry): %s", queueName, string(body))
			return nil
		}
		log.Println("❌ Message not acknowledged by broker (after retry)")
		return amqp.ErrClosed
	case <-time.After(5 * time.Second):
		log.Println("⚠️ Timeout waiting for broker confirmation (after retry). Treating as success.")
		log.Printf("📤 Sent (no confirm) message to %s (after retry): %s", queueName, string(body))
		return nil
	}
}

// GetChannel returns the active AMQP channel (may be nil)
func GetChannel() *amqp.Channel {
	mutex.Lock()
	defer mutex.Unlock()
	return channel
}

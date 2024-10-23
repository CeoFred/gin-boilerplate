package streaming

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type Config struct {
	Version   string
	Brokers   string
	Topic     []string
	Producers int
	Verbose   bool
	// consumer config
	Assignor string
	Oldest   bool
	Group    string
	Ctx      context.Context
}

type EventProducer interface {
	BroadCast(count int, eventName string, payload []byte) error
	Clear()
}

func NewProducer(c *Config) (EventProducer, error) {
	p := &Producer{}
	p.setUp(c)
	return p, nil
}

type Producer struct {
	producerProvider *producerProvider
}

func (p *Producer) BroadCast(count int, eventName string, payload []byte) error {
	producer := p.producerProvider.borrow()
	defer p.producerProvider.release(producer)

	err := producer.BeginTxn()
	if err != nil {
		log.Printf("unable to start txn %s\n", err)
		return err
	}

	producer.Input() <- &sarama.ProducerMessage{Topic: eventName, Key: nil, Value: sarama.StringEncoder(payload)}

	err = producer.CommitTxn()
	if err != nil {
		log.Printf("Producer: unable to commit txn %s\n", err)
		for {
			if producer.TxnStatus()&sarama.ProducerTxnFlagFatalError != 0 {
				// fatal error. need to recreate producer.
				log.Printf("Producer: producer is in a fatal state, need to recreate it")
				return err

			}
			// If producer is in abortable state, try to abort current transaction.
			if producer.TxnStatus()&sarama.ProducerTxnFlagAbortableError != 0 {
				err = producer.AbortTxn()
				if err != nil {
					// If an error occured just retry it.
					log.Printf("Producer: unable to abort transaction: %+v", err)

					return err

				}
				break
			}
			// if not you can retry
			err = producer.CommitTxn()
			if err != nil {
				log.Printf("Producer: unable to commit txn %s\n", err)
				return err
			}
		}
		return nil
	}
	return nil
}

func (p *Producer) Clear() {
	p.producerProvider.clear()
}

func (p *Producer) setUp(c *Config) {
	log.Println("Starting a new Sarama producer")

	if c.Verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	version, err := sarama.ParseKafkaVersion(c.Version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	producerProvider := newProducerProvider(strings.Split(c.Brokers, ","), func() *sarama.Config {
		config := sarama.NewConfig()
		config.Version = version
		config.Producer.Idempotent = true
		config.Producer.Return.Errors = false
		config.Producer.RequiredAcks = sarama.WaitForAll
		config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
		config.Producer.Transaction.Retry.Backoff = 10
		config.Producer.Transaction.ID = "txn_producer"
		config.Net.MaxOpenRequests = 1
		config.Net.DialTimeout = time.Second * 10
		config.ClientID = "gin-app"
		return config
	})
	p.producerProvider = producerProvider
}

type producerProvider struct {
	transactionIdGenerator int32

	producersLock sync.Mutex
	producers     []sarama.AsyncProducer

	producerProvider func() sarama.AsyncProducer
}

func newProducerProvider(brokers []string, producerConfigurationProvider func() *sarama.Config) *producerProvider {
	provider := &producerProvider{}
	provider.producerProvider = func() sarama.AsyncProducer {
		config := producerConfigurationProvider()
		suffix := provider.transactionIdGenerator
		// Append transactionIdGenerator to current config.Producer.Transaction.ID to ensure transaction-id uniqueness.
		if config.Producer.Transaction.ID != "" {
			provider.transactionIdGenerator++
			config.Producer.Transaction.ID = config.Producer.Transaction.ID + "-" + fmt.Sprint(suffix)
		}
		producer, err := sarama.NewAsyncProducer(brokers, config)
		if err != nil {
			return nil
		}
		return producer
	}
	return provider
}

func (p *producerProvider) borrow() (producer sarama.AsyncProducer) {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	if len(p.producers) == 0 {
		for {
			producer = p.producerProvider()
			if producer != nil {
				return
			}
		}
	}

	index := len(p.producers) - 1
	producer = p.producers[index]
	p.producers = p.producers[:index]
	return
}

func (p *producerProvider) release(producer sarama.AsyncProducer) {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	// If released producer is erroneous close it and don't return it to the producer pool.
	if producer.TxnStatus()&sarama.ProducerTxnFlagInError != 0 {
		// Try to close it
		err := producer.Close()
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	p.producers = append(p.producers, producer)
}

func (p *producerProvider) clear() {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	for _, producer := range p.producers {
		err := producer.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
	p.producers = p.producers[:0]
}

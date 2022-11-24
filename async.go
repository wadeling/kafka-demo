package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/wadeling/kafka-demo/pkg/msg"
	"log"
	"sync"
	"time"
)

func initAsyncProducer(brokerList []string) (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Flush.Frequency = 500 * time.Millisecond
	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Printf("failed to create asyncProducer:%v\n", err)
		return nil, err
	}
	return producer, err
}

func asyncSendMsg(brokerList []string) {
	wg := sync.WaitGroup{}
	nodeNum := 2
	imageNum := 2
	for i := 0; i < nodeNum; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			producer, err := initAsyncProducer(brokerList)
			if err != nil {
				return
			}

			tmsg := msg.GetMockMsg()
			log.Printf("node %d start send kafka msg.", index)
			for j := 0; j < imageNum; j++ {
				producer.Input() <- &sarama.ProducerMessage{
					Topic: Topic,
					Key:   sarama.StringEncoder(fmt.Sprintf("wade-%d", index)),
					Value: &tmsg,
				}
			}

			subwg := sync.WaitGroup{}
			subwg.Add(1)
			go func() {
				defer subwg.Done()
				for err := range producer.Errors() {
					log.Println("Failed to write entry:", err)
				}
			}()
			subwg.Wait()
			log.Printf("node %d end send kafka msg.", index)
		}(i)
	}

	wg.Wait()

}

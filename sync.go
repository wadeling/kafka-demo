package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"sync"
)

func initSyncProducer(brokerList []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
		return nil, err
	}
	return producer, nil
}

func syncSendMsg(brokerList []string) {
	log.Printf("sync send,brokers:%v\n", brokerList)

	wg := sync.WaitGroup{}
	nodeNum := 2
	imageNum := 2
	for i := 0; i < nodeNum; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			producer, err := initSyncProducer(brokerList)
			if err != nil {
				log.Printf("init sync producer err:%v", err)
				return
			}

			msg := getMockMsg()
			log.Printf("node %d start send kafka msg.", index)
			for j := 0; j < imageNum; j++ {
				partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
					Topic: Topic,
					Key:   sarama.StringEncoder(fmt.Sprintf("wade-sync-%d", index)),
					Value: &msg,
				})
				if err != nil {
					log.Printf("failed to send sync msg:%v", err)
					continue
				}
				log.Printf("node %d send msg %d ok,partitions %v,offset %v\n", index, j, partition, offset)
			}

			log.Printf("node %d end send kafka msg.", index)
		}(i)
	}

	wg.Wait()

}

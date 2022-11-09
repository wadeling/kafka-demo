package main

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"strings"
	"time"
)

const (
	Topic = "wade-test"
)

var (
	Brokers = []string{"localhost:9092"}
)

// NodeImage test msg struct
type NodeImage struct {
	Repo      string
	Tag       string
	Digest    string
	CreatedAt int64
	encoded   []byte
	err       error
}

func (n *NodeImage) ensureEncoded() {
	if n.err == nil && n.encoded == nil {
		n.encoded, n.err = json.Marshal(n)
	}
}
func (n *NodeImage) Length() int {
	n.ensureEncoded()
	return len(n.encoded)
}
func (n *NodeImage) Encode() ([]byte, error) {
	n.ensureEncoded()
	return n.encoded, n.err
}

func getMockMsg() NodeImage {
	nodeImage := NodeImage{
		Repo:      "index.docker.io/wade23",
		Tag:       "test",
		Digest:    "45a7d023456f582161e545cac6a07d9052d481cfb87dcc2b5c81b6e052dfd5d5",
		CreatedAt: time.Now().UnixMicro(),
	}
	return nodeImage
}

func createTopic(brokerList []string) error {
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	admin, err := sarama.NewClusterAdmin(brokerList, config)
	if err != nil {
		log.Fatal("Error while creating cluster admin: ", err.Error())
		return err
	}
	defer func() { _ = admin.Close() }()
	err = admin.CreateTopic(Topic, &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}, false)
	if err != nil {
		log.Fatal("Error while creating topic: ", err.Error())
		return err
	}
	return nil
}

func main() {
	log.Print("start")

	if len(os.Args) < 2 {
		log.Printf("Usage:%s {create|sync-send|async-send} [brokers]", os.Args[0])
		return
	}

	// get broker list from cmd
	var brokerList []string
	if len(os.Args) >= 3 {
		// get brokers
		arr := strings.Split(os.Args[2], ",")
		brokerList = append(brokerList, arr...)
	} else {
		// use default addr
		brokerList = Brokers
	}
	log.Printf("broker list:%v\n", brokerList)

	if os.Args[1] == "create" {
		err := createTopic(brokerList)
		if err != nil {
			log.Printf("create topic err:%v\n", err)
			return
		}
		log.Print("create topic ok.\n")
	} else if os.Args[1] == "async-send" {
		asyncSendMsg(brokerList)
	} else if os.Args[1] == "sync-send" {
		syncSendMsg(brokerList)
	} else if os.Args[1] == "consume" {
		//todo: consume data
	} else {
		log.Printf("not support action:%s", os.Args[1])
	}

	log.Print("end")
}

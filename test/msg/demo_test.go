package demo

import (
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/wadeling/kafka-demo/pkg/msg"
	"testing"
	"time"
)

var (
	repoName = "index.docker.io/wade23"
	tag      = "test"
	digest   = "45a7d023456f582161e545cac6a07d9052d481cfb87dcc2b5c81b6e052dfd5d5"
)

type CustomMsg struct {
	Key   []byte
	Value []byte
}

func TestMsgSize(t *testing.T) {
	t.Log("start")
	now := time.Now().UnixMicro()
	tmpData := ""
	for i := 0; i < 40*1000; i++ {
		tmpData = tmpData + "a"
	}
	nodeImage := msg.NodeImage{
		Repo:      repoName,
		Tag:       tag,
		Digest:    digest,
		CreatedAt: now,
		Data:      tmpData,
	}
	orgLen := len(repoName) + len(tag) + len(digest) + len(tmpData) + 8
	t.Logf("orglen %v", orgLen)

	data, err := json.Marshal(nodeImage)
	if err != nil {
		t.Fatalf("marchal failed.%v", err)
	}
	t.Logf("data len:%v", len(data))
	tmsg := kafka.Message{
		Key:   []byte("aaa"),
		Value: data,
	}

	msgData, err := json.Marshal(tmsg)
	if err != nil {
		t.Fatalf("marcharl failed.%v", err)
	}
	t.Logf("msg data len:%v", len(msgData))

	cmsg := CustomMsg{
		Key:   []byte("aaa"),
		Value: data,
	}
	msgData2, err := json.Marshal(cmsg)
	if err != nil {
		t.Fatalf("marshal failed %v", err)
	}
	t.Logf("cmsg data len:%v", len(msgData2))
}

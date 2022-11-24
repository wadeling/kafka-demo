package msg

import (
	"encoding/json"
	"time"
)

// NodeImage test msg struct
type NodeImage struct {
	Repo      string
	Tag       string
	Digest    string
	CreatedAt int64
	Data      string
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

func GetMockMsg() NodeImage {
	// mock data
	data := ""
	for i := 0; i < 32*1000; i++ {
		data = data + "a"
	}
	nodeImage := NodeImage{
		Repo:      "index.docker.io/wade23",
		Tag:       "test",
		Digest:    "45a7d023456f582161e545cac6a07d9052d481cfb87dcc2b5c81b6e052dfd5d5",
		CreatedAt: time.Now().UnixMicro(),
		Data:      data,
	}
	return nodeImage
}

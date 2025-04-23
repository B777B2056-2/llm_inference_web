package client

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"llm_online_inference/scheduler/confparser"
	"llm_online_inference/scheduler/resource"
	"sync"
	"time"
)

var responseChanMap sync.Map

const (
	ChatCompletionKafkaHeaderKey = "ChatSessionID"
	TopicPrompt                  = "user-prompt"
	TopicAnswer                  = "model-answer"
)

type KafkaReadWriteClient struct {
	addresses    []string
	responseChan chan string
}

func NewKafkaReadWriteClient() *KafkaReadWriteClient {
	var addresses []string
	for _, item := range confparser.ResourceConfig.Kafka.Addresses {
		addresses = append(addresses, fmt.Sprintf("%s:%d", item.Host, item.Port))
	}
	return &KafkaReadWriteClient{
		addresses:    addresses,
		responseChan: make(chan string, 100),
	}
}

// Send 推送一次消息
func (kw *KafkaReadWriteClient) Send(ctx context.Context, userID int, topic string,
	header kafka.Header, key, data string) error {
	writeTimeout := time.Duration(confparser.ResourceConfig.Kafka.SendMsgTimeoutInMs) * time.Millisecond
	writer := &kafka.Writer{
		Addr:         kafka.TCP(kw.addresses...),
		Topic:        topic,
		Balancer:     &kafka.CRC32Balancer{},
		Compression:  kafka.Gzip,
		WriteTimeout: time.Duration(confparser.ResourceConfig.Kafka.SendMsgTimeoutInMs) * time.Millisecond,
	}

	msg := kafka.Message{
		Headers: []kafka.Header{header},
		Key:     []byte(key),
		Value:   []byte(data),
	}

	ctx, cancel := context.WithTimeout(ctx, writeTimeout)
	defer cancel()

	if err := writer.WriteMessages(ctx, msg); err != nil {
		resource.Logger.WithFields(
			logrus.Fields{
				"userID": userID,
				"topic":  topic,
				"err":    err,
			},
		).Error("failed to write messages to kafka")
		return err
	}
	if err := writer.Close(); err != nil {
		resource.Logger.WithFields(
			logrus.Fields{
				"userID": userID,
				"topic":  topic,
				"err":    err,
			},
		).Error("failed to close kafka writer")
		return err
	}

	// 建立写时header与读时channel的关系
	responseChanMap.Store(header.Key, kw.responseChan)
	return nil
}

func getHeader(headers []kafka.Header, key string) string {
	for _, h := range headers {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}

// Read 读取消息
func (kw *KafkaReadWriteClient) Read(ctx context.Context, userId int, topic, key, groupID string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  kw.addresses,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 1,
		MaxBytes: 10e6, // 10MB
	})
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			resource.Logger.WithFields(
				logrus.Fields{
					"userId": userId,
					"topic":  topic,
					"err":    err,
				},
			).Error("failed to read messages from kafka")
			continue
		}

		// 处理消息：塞入channel，并提交偏移量（不匹配时不提交偏移量，使得下次该消息仍然能被消费）
		correlationID := getHeader(msg.Headers, key)
		if ch, ok := responseChanMap.Load(correlationID); ok {
			ch.(chan string) <- string(msg.Value)

			if err := reader.CommitMessages(ctx, msg); err != nil {
				resource.Logger.WithFields(
					logrus.Fields{
						"userId": userId,
						"topic":  topic,
						"msg":    msg,
						"err":    err,
					},
				).Error("failed to commit offset to kafka")
				continue
			}
		}
	}
}

// GetReadChan 获取异步消息流读取channel
func (kw *KafkaReadWriteClient) GetReadChan() <-chan string {
	return kw.responseChan
}

// ReadCompletedHandler 读取消息完成后调用
func (kw *KafkaReadWriteClient) ReadCompletedHandler(key string) {
	responseChanMap.Delete(key)
	close(kw.responseChan)
}

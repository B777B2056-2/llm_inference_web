package client

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"llm_online_inference/accessor/confparser"
	"llm_online_inference/accessor/resource"
	"time"
)

const (
	TopicBatchInferenceRequests = "batch-inference-requests"
)

type KafkaProductorClient struct {
	topic        string
	writer       *kafka.Writer
	writeTimeout time.Duration
}

func NewKafkaProductorClient(topic string) *KafkaProductorClient {
	var addresses []string
	for _, item := range confparser.ResourceConfig.Kafka.Addresses {
		addresses = append(addresses, fmt.Sprintf("%s:%d", item.Host, item.Port))
	}
	writeTimeout := time.Duration(confparser.ResourceConfig.Kafka.SendMsgTimeoutInMs) * time.Millisecond
	writer := &kafka.Writer{
		Addr:         kafka.TCP(addresses...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		Compression:  kafka.Gzip,
		WriteTimeout: writeTimeout,
	}
	return &KafkaProductorClient{topic: topic, writer: writer, writeTimeout: writeTimeout}
}

// Send 推送一次消息
func (kw *KafkaProductorClient) Send(ctx context.Context, msgs []kafka.Message) error {
	ctx, cancel := context.WithTimeout(ctx, kw.writeTimeout)
	defer cancel()

	if err := kw.writer.WriteMessages(ctx, msgs...); err != nil {
		resource.Logger.WithFields(
			logrus.Fields{
				"topic": kw.topic,
				"err":   err,
			},
		).Error("failed to write messages to kafka")
		return err
	}
	if err := kw.writer.Close(); err != nil {
		resource.Logger.WithFields(
			logrus.Fields{
				"topic": kw.topic,
				"err":   err,
			},
		).Error("failed to close kafka writer")
		return err
	}
	return nil
}

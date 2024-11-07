package main

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	logger "load-test-lab/pkg"
	"log"
	"os"
	"time"
)

func main() {
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	kafkaBootstrap := os.Getenv("KAFKA_BOOTSTRAP")
	consumerGroupName := os.Getenv("CONSUMER_GROUP_NAME")
	mongoURI := os.Getenv("MONGO_URI")
	mongoDatabase := os.Getenv("MONGO_DATABASE")
	mongoCollection := os.Getenv("MONGO_COLLECTION")
	mongoUser := os.Getenv("MONGO_USER")
	mongoPass := os.Getenv("MONGO_PASS")

	var errLog = logger.NewErrorFile("kafka_consumer")

	fmt.Printf("KAFKA_TOPIC: %v\n"+
		"KAFKA_BOOTSTRAP: %v\n"+
		"CONSUMER_GROUP_NAME: %v\n"+
		"MONGO_URI: %v\n"+
		"MONGO_DATABASE: %v\n"+
		"MONGO_COLLECTION: %v\n", kafkaTopic, kafkaBootstrap, consumerGroupName, mongoURI, mongoDatabase, mongoCollection)

	ctx := context.Background()

	var cred options.Credential
	cred.Username = mongoUser
	cred.Password = mongoPass

	clientOptions := options.Client().ApplyURI(mongoURI).SetAuth(cred)
	mongoClient, err := mongo.Connect(clientOptions)
	if err != nil {
		errLog.Fatalf("Failed to create mongo client: %v", err)
	}
	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			errLog.Fatalf("Failed to disconnect mongo client: %v", err)
		}
	}()

	collection := mongoClient.Database(mongoDatabase).Collection(mongoCollection)

	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	consumerGroup, err := sarama.NewConsumerGroup([]string{kafkaBootstrap}, consumerGroupName, config)
	if err != nil {
		errLog.Fatalf("Failed to create consumer group: %v", err)
	}
	defer func() {
		if err = consumerGroup.Close(); err != nil {
			errLog.Fatalf("Failed to close consumer group: %v", err)
		}
	}()

	// Consumer handler
	go func() {
		for {
			if err := consumerGroup.Consume(context.Background(), []string{kafkaTopic}, &consumer{collection: collection}); err != nil {
				errLog.Fatalf("Failed to start consuming from kafka: %v", err)
			}
		}
	}()

	// Wait forever
	for {
	}
}

var _ sarama.ConsumerGroupHandler = (*consumer)(nil)

// consumer implements the sarama.ConsumerGroupHandler interface
type consumer struct {
	collection *mongo.Collection
}

func (c *consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// Форматируем данные для записи в MongoDB
		document := bson.M{
			"key":   string(msg.Key),
			"value": string(msg.Value),
			"ts":    time.Now(),
		}

		_, err := c.collection.InsertOne(context.TODO(), document)
		if err != nil {
			log.Printf("Ошибка при записи в MongoDB: %v\n", err)
		}
		session.MarkMessage(msg, "")
	}
	return nil
}

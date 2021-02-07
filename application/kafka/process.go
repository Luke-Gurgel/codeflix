package kafka

import (
	"fmt"
	"os"

	"github.com/Luke-Gurgel/codeflix/application/factory"
	appmodel "github.com/Luke-Gurgel/codeflix/application/model"
	"github.com/Luke-Gurgel/codeflix/application/usecase"
	"github.com/Luke-Gurgel/codeflix/domain/model"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jinzhu/gorm"
)

type KafkaProcessor struct {
	Database        *gorm.DB
	Producer        *ckafka.Producer
	DeliveryChannel chan ckafka.Event
}

func CreateKafkaProcessor(database *gorm.DB, producer *ckafka.Producer, deliveryChannel chan ckafka.Event) *KafkaProcessor {
	return &KafkaProcessor{
		Database:        database,
		Producer:        producer,
		DeliveryChannel: deliveryChannel,
	}
}

func (k *KafkaProcessor) Consume() {
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": os.Getenv("kafkaBootstrapServers"),
		"group.id":          os.Getenv("kafkaConsumerGroupID"),
		"auto.offset.reset": "earliest",
	}

	consumer, err := ckafka.NewConsumer(configMap)

	if err != nil {
		panic(err)
	}

	topics := []string{os.Getenv("kafkaTransactionTopic"), os.Getenv("kafkaTransactionConfirmationTopic")}
	err = consumer.SubscribeTopics(topics, nil)

	if err == nil {
		fmt.Println("kafka consumer has been started")
	} else {
		panic(err)
	}

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			fmt.Println(string(msg.Value))
			k.processMessage(msg)
		} else {
			panic(err)
		}
	}
}

func (k *KafkaProcessor) processMessage(msg *ckafka.Message) {
	transactionsTopic := "transactions"
	transactionsConfirmationTopic := "transaction_confirmation"

	switch topic := *msg.TopicPartition.Topic; topic {
	case transactionsTopic:
		k.processTransaction(msg)
	case transactionsConfirmationTopic:
		k.processTransactionConfirmation(msg)
	default:
		fmt.Println("Not a valid topic", string(msg.Value))
	}
}

func (k *KafkaProcessor) processTransaction(msg *ckafka.Message) error {
	transaction := appmodel.CreateTransaction()
	err := transaction.ParseJson(msg.Value)

	if err != nil {
		return err
	}

	transactionUseCase := factory.TransactionUseCase(k.Database)
	createdTransaction, err := transactionUseCase.Register(
		transaction.AccountID,
		transaction.Amount,
		transaction.PixKeyTo,
		transaction.PixKeyToKind,
		transaction.Description,
	)

	if err != nil {
		fmt.Println("Error registering transactiond", err)
		return err
	}

	topic := "bank" + createdTransaction.PixKeyTo.Account.Bank.Code
	transaction.ID = createdTransaction.ID
	transaction.Status = model.TransactionPending
	transactionJson, err := transaction.ToJson()

	if err != nil {
		return err
	}

	err = Publish(string(transactionJson), topic, k.Producer, k.DeliveryChannel)

	if err != nil {
		return err
	}

	return nil
}

func (k *KafkaProcessor) processTransactionConfirmation(msg *ckafka.Message) error {
	transaction := appmodel.CreateTransaction()
	err := transaction.ParseJson(msg.Value)

	if err != nil {
		return err
	}

	transactionUseCase := factory.TransactionUseCase(k.Database)

	if transaction.Status == model.TransactionConfirmed {
		err = k.confirmTransaction(transaction, transactionUseCase)
		if err != nil {
			return err
		}
	} else if transaction.Status == model.TransactionCompleted {
		_, err := transactionUseCase.Complete(transaction.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k *KafkaProcessor) confirmTransaction(transaction *appmodel.Transaction, transactionUseCase usecase.TransactionUseCase) error {
	confirmedTransaction, err := transactionUseCase.Confirm(transaction.ID)

	if err != nil {
		return err
	}

	topic := "bank" + confirmedTransaction.AccountFrom.Bank.Code
	transactionJson, err := transaction.ToJson()

	if err != nil {
		return err
	}

	err = Publish(string(transactionJson), topic, k.Producer, k.DeliveryChannel)

	if err != nil {
		return err
	}

	return nil
}

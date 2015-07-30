package mq

import "github.com/streadway/amqp"

type RabbitChannel struct {
	inner *amqp.Channel
}

func newChannel(inner *amqp.Channel) *RabbitChannel {
	return &RabbitChannel{inner: inner}
}

func (this *RabbitChannel) ConfigureChannelBuffer(messageCount int) error {
	return this.inner.Qos(messageCount, 0, false)
}
func (this *RabbitChannel) ConfigureChannelAsTransactional() error {
	return this.inner.Tx()
}

func (this *RabbitChannel) DeclareTransientQueue() (string, error) {
	if queue, err := this.inner.QueueDeclare("", false, true, false, false, nil); err != nil {
		return "", err
	} else {
		return queue.Name, nil
	}
}
func (this *RabbitChannel) BindExchangeToQueue(queue, exchange string) error {
	return this.inner.QueueBind(queue, "", exchange, false, nil)
}

func (this *RabbitChannel) Consume(queueName, consumerName string) (<-chan amqp.Delivery, error) {
	return this.inner.Consume(queueName, consumerName, false, false, false, false, nil)
}
func (this *RabbitChannel) ConsumeWithoutAcknowledgement(queueName, consumerName string) (<-chan amqp.Delivery, error) {
	return this.inner.Consume(queueName, consumerName, true, true, false, false, nil)
}
func (this *RabbitChannel) ExclusiveConsume(queueName, consumerName string) (<-chan amqp.Delivery, error) {
	return this.inner.Consume(queueName, consumerName, false, true, false, false, nil)
}
func (this *RabbitChannel) ExclusiveConsumeWithoutAcknowledgement(queueName, consumerName string) (<-chan amqp.Delivery, error) {
	return this.inner.Consume(queueName, consumerName, true, true, false, false, nil)
}
func (this *RabbitChannel) CancelConsumer(consumerName string) error {
	return this.inner.Cancel(consumerName, false)
}

func (this *RabbitChannel) AcknowledgeSingleMessage(deliveryTag uint64) error {
	return this.inner.Ack(deliveryTag, false)
}
func (this *RabbitChannel) AcknowledgeMultipleMessages(deliveryTag uint64) error {
	return this.inner.Ack(deliveryTag, true)
}

func (this *RabbitChannel) PublishMessage(destination string, message amqp.Publishing) error {
	return this.inner.Publish(destination, "", false, false, message)
}

func (this *RabbitChannel) CommitTransaction() error {
	return this.inner.TxCommit()
}
func (this *RabbitChannel) RollbackTransaction() error {
	return this.inner.TxRollback()
}

func (this *RabbitChannel) Close() error {
	return this.inner.Close()
}

package messaging

import (
	"reflect"
	"strings"
)

type DispatchWriter struct {
	writer    Writer
	committer CommitWriter
	discovery TypeDiscovery
	overrides map[reflect.Type]Dispatch
	template  Dispatch
}

func NewDispatchWriter(writer Writer, discovery TypeDiscovery) *DispatchWriter {
	committer, _ := writer.(CommitWriter)

	return &DispatchWriter{
		writer:    writer,
		committer: committer,
		discovery: discovery,
		overrides: make(map[reflect.Type]Dispatch),
		template:  Dispatch{Durable: true},
	}
}

func (this *DispatchWriter) RegisterTemplate(template Dispatch) {
	this.template = template
}

func (this *DispatchWriter) RegisterOverride(instanceType reflect.Type, message Dispatch) {
	this.overrides[instanceType] = message
}

func (this *DispatchWriter) Write(item Dispatch) error {
	messageType := reflect.TypeOf(item.Message)

	target, found := this.overrides[messageType]
	if !found {
		target = this.template
		if discovered, err := this.discovery.Discover(item.Message); err != nil {
			return err
		} else {
			target.MessageType = discovered
		}

		target.Destination = strings.Replace(target.MessageType, ".", "-", -1)
	}

	target.Message = item.Message
	return this.writer.Write(target)
}

func (this *DispatchWriter) Commit() error {
	if this.committer == nil {
		return nil
	}

	return this.committer.Commit()
}

func (this *DispatchWriter) Close() {
	this.writer.Close()
}

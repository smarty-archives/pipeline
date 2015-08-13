package messaging

import (
	"reflect"
	"strings"
)

type SerializationWriter struct {
	writer            Writer
	commitWriter      CommitWriter
	serializer        Serializer
	messageTypePrefix string
}

func NewSerializationWriter(inner Writer, serializer Serializer, messageTypePrefix string) *SerializationWriter {
	commitWriter, _ := inner.(CommitWriter)
	return &SerializationWriter{
		writer:            inner,
		commitWriter:      commitWriter,
		serializer:        serializer,
		messageTypePrefix: messageTypePrefix,
	}
}

func (this *SerializationWriter) Write(dispatch Dispatch) error {
	if len(dispatch.Payload) > 0 && len(dispatch.MessageType) > 0 {
		return this.writer.Write(dispatch) // already have a payload a message type, forward to inner
	}

	if len(dispatch.Payload) == 0 && dispatch.Message == nil {
		return EmptyDispatchError // no payload and no message, this is a total fail
	}

	if dispatch.Message == nil && len(dispatch.MessageType) == 0 {
		return MessageTypeDiscoveryError // no message type and no way to get it
	}

	if payload, err := this.serializer.Serialize(dispatch.Message); err != nil {
		return err // serialization failed
	} else {
		dispatch.Payload = payload
	}

	if len(dispatch.MessageType) == 0 {
		dispatch.MessageType = discoverType(this.messageTypePrefix, dispatch.Message)
	}

	// TODO
	// if len(dispatch.MessageType) == 0 {
	// 	return errors.New("Unable to write message, the type cannot be discovered.")
	// }

	return this.writer.Write(dispatch)
}

// TODO: should this be pulled out into its own stucture?
// that would allow different behaviors to be invoked
func discoverType(prefix string, message interface{}) string {
	reflectType := reflect.TypeOf(message)
	if name := reflectType.Name(); len(name) > 0 {
		return prefix + strings.ToLower(name)
	}

	name := reflectType.String()
	index := strings.LastIndex(name, ".")

	if index == -1 {
		return ""
	}

	suffix := strings.ToLower(name[index+1:])
	if strings.HasPrefix(name, "*") {
		return "*" + prefix + suffix
	} else {
		return prefix + suffix
	}
}

func (this *SerializationWriter) Commit() error {
	if this.commitWriter == nil {
		return nil
	}

	return this.commitWriter.Commit()
}

func (this *SerializationWriter) Close() {
	this.writer.Close()
}

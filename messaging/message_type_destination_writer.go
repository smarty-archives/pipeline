package messaging

import "strings"

// TODO: this structure is not under test at all
type MessageTypeDestinationWriter struct {
	writer       Writer
	commitWriter CommitWriter
}

func NewMessageTypeDestinationWriter(inner Writer) *MessageTypeDestinationWriter {
	commitWriter, _ := inner.(CommitWriter)
	return &MessageTypeDestinationWriter{writer: inner, commitWriter: commitWriter}
}

func (this *MessageTypeDestinationWriter) Write(dispatch Dispatch) error {
	if len(dispatch.Destination) > 0 {
		return this.writer.Write(dispatch)
	}

	if len(dispatch.MessageType) == 0 {
		return UnroutableDispatchError
	}

	dispatch.Destination = strings.Replace(dispatch.MessageType, ".", "-", -1)
	dispatch.Destination = strings.ToLower(dispatch.Destination)
	return this.writer.Write(dispatch)
}

func (this *MessageTypeDestinationWriter) Commit() error {
	if this.commitWriter == nil {
		return nil
	}

	return this.commitWriter.Commit()
}

func (this *MessageTypeDestinationWriter) Close() {
	this.writer.Close()
}

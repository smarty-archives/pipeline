package messaging

type SerializationWriter struct {
	serializer Serializer
}

func NewSerializationWriter(serializer Serializer) *SerializationWriter {
	return &SerializationWriter{
		serializer: serializer,
	}
}

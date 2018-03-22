package channelling

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"

	"buffercache"
)

type IncomingDecoder interface {
	DecodeIncoming(buffercache.Buffer) (*DataIncoming, error)
}

type OutgoingEncoder interface {
	EncodeOutgoing(*DataOutgoing) (buffercache.Buffer, error)
}

type Codec interface {
	NewBuffer() buffercache.Buffer
	IncomingDecoder
	OutgoingEncoder
}

type incomingCodec struct {
	buffers       buffercache.BufferCache
	incomingLimit int
}

func NewCodec(incomingLimit int) Codec {
	return &incomingCodec{buffercache.NewBufferCache(1024, bytes.MinRead), incomingLimit}
}

func (codec incomingCodec) NewBuffer() buffercache.Buffer {
	return codec.buffers.New()
}

func (codec incomingCodec) DecodeIncoming(b buffercache.Buffer) (*DataIncoming, error) {
	length := b.GetBuffer().Len()
	if length > codec.incomingLimit {
		return nil, errors.New("Incoming message size limit exceeded")
	}
	incoming := &DataIncoming{}
	return incoming, json.Unmarshal(b.Bytes(), incoming)
}

func (codec incomingCodec) EncodeOutgoing(outgoing *DataOutgoing) (buffercache.Buffer, error) {
	b := codec.NewBuffer()
	if err := json.NewEncoder(b).Encode(outgoing); err != nil {
		log.Println("Error while encoding JSON", err)
		b.Decref()
		return nil, err
	}
	return b, nil
}

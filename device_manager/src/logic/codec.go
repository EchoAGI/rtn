package logic



type Encoder interface {
	Encode(msg interface{}) []byte
}

type Decoder interface {
	Decode(bytes []byte, out interface{})
}

type Codec interface {
	Encoder
	Decoder
}
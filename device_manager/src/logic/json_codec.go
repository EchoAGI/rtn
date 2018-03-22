package logic


type JsonCodec struct {

}

func NewJsonCodec() *JsonCodec {
	return &JsonCodec{

	}
}


func (c *JsonCodec) Encode(msg interface{}) []byte {
	return []byte{}
}

func (c *JsonCodec) Decode(bytes []byte, out interface{}) {
	out = "test"
}
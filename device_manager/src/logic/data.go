package logic



type DataRegisterDevice struct {
	DeviceId string
	Topic string
	Extra    map[string]interface{}
}

type DataUnRegisterDevice struct {
	DeviceId string
	Extra    map[string]interface{}
}

type DataIncoming struct {
	Token string
	ClientId string
	Extra    map[string]interface{}
}


type DataMessage struct {
	Type int
	RegisterDevice DataRegisterDevice
	UnRegisterDevice DataUnRegisterDevice
	DataIncoming DataIncoming
}
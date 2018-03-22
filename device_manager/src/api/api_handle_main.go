package api

import(
	"logic"
)


type APIImpl struct {

}

func NewAPIImpl() *APIImpl {
	return &APIImpl{}
}


func(api *APIImpl) HandleRegisterDevice(client *logic.Client, msg *logic.DataRegisterDevice) {

}
func (api *APIImpl) HandleUnRegisterDevice(client *logic.Client, msg *logic.DataUnRegisterDevice) {

}
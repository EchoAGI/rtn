package api

import(
	"logic"
)


type API interface {
	HandleRegisterDevice(client *logic.Client, msg *logic.DataRegisterDevice)
	HandleUnRegisterDevice(client *logic.Client, msg *logic.DataUnRegisterDevice)
}
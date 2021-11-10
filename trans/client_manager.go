package trans

type ClientManager struct {
	DeviceManager
}

var clientManagerInstance *DeviceManager
func ClientManagerInstance() *DeviceManager{
	if clientManagerInstance == nil {
		clientManagerInstance = DeviceManagerNew()
	}
	return clientManagerInstance
}

func ClientManagerNew()* ClientManager{
	return &ClientManager{}
}
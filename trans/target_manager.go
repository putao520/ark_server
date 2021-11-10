package trans

type TargetManager struct {
	DeviceManager
}

var targetManagerInstance *DeviceManager
func TargetManagerInstance() *DeviceManager{
	if targetManagerInstance == nil {
		targetManagerInstance = DeviceManagerNew()
	}
	return targetManagerInstance
}

func TargetManagerNew()* TargetManager{
	return &TargetManager{}
}
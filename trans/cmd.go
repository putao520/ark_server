package trans

type resultCmd string

const (
	DeviceConnected  resultCmd = "connect"
	DeviceDisconnect resultCmd = "disconnect"

	// InvokeString 运行JS命令
	InvokeString resultCmd = "invokeString"
	// InvokeScript 运行 预制脚本
	InvokeScript resultCmd = "invokeScript"

	// HealthInfo 当前任务状态信息
	HealthInfo resultCmd = "healthInfo"

	// TargetReconnect 目标设备重新上线
	TargetReconnect resultCmd = "targetReconnect"
)

type ResultCode int

const (
	Success ResultCode = 0
	Failed  ResultCode = 1
)
package trans

import (
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

func checkDeviceQueue(queue map[string]*DeviceObject, delay int64) {
	pingContent := ""
	for true {
		time.Sleep(time.Duration(delay) * time.Second)
		// queue :=TargetManagerInstance().targetQueue
		for deviceId, do := range queue {
			// err := do.client.WriteControl(websocket.PingMessage, []byte(pingContent), time.Now().Add(10 *time.Second))
			err := do.client.WriteMessage(websocket.TextMessage, []byte(pingContent))
			if err != nil {
				fmt.Println(err)
				do.LeaveSession()
				delete(queue, deviceId)
			}
			fmt.Println(deviceId)
		}
	}
}

var checked = false

func StartChecker() {
	if checked == false {
		checked = true
		go checkDeviceQueue(TargetManagerInstance().targetQueue, 180)
		go checkDeviceQueue(ClientManagerInstance().targetQueue, 3600)
	}
}

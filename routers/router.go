// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"server/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// Http.
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/scriptLibrary",
			beego.NSInclude(
					&controllers.ScriptLibraryController{},
				),
			),
	)
	beego.AddNamespace(ns)
	// http
	beego.Router("/room", &controllers.RoomController{})

	// WebSocket.
	beego.Router("/ws", &controllers.WebSocketController{}, "get:OnConnect")
	// beego.Router("/ws/Register", &controllers.WebSocketController{}, "get:onConnect")
}

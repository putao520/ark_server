// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"server/controllers"
	"server/filter"
	"server/jwt"
)

func init() {

	beego.InsertFilter("/v1/room/*", beego.BeforeExec, jwt.FilterJwt)
	beego.InsertFilter("/v1/script/*", beego.BeforeExec, jwt.FilterJwt)
	beego.InsertFilter("/v1/user/*", beego.BeforeExec, filter.UserFilter)

	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/script",
			beego.NSInclude(
				&controllers.ScriptController{},
			),
		),

		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),

		beego.NSNamespace("/room",
			beego.NSInclude(
				&controllers.RoomController{},
			),
		),
	)
	beego.AddNamespace(ns)

	// WebSocket.
	beego.Router("/ws", &controllers.WebSocketController{}, "get:OnConnect")
}

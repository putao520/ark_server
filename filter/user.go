package filter

import (
	"github.com/beego/beego/v2/server/web/context"
	"server/jwt"
)

func UserFilter(ctx *context.Context) {
	if ctx.Request.Method == "PUT" {
		jwt.FilterJwt(ctx)
	}
}

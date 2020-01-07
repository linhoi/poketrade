package middleware

import (
	"github.com/kataras/iris/v12"
)

func AuthBeforeLogin(ctx iris.Context) {
	uid := ctx.GetCookie("uid")
	if uid == ""{
		ctx.Application().Logger().Debug("You Must Login First")
		ctx.Redirect("/user/login")
	}
	ctx.Application().Logger().Debug("Login Success")
	ctx.Next()
}

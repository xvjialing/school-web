// @APIVersion 1.0.0
// @Title school API
// @Description school API
// @Contact xvjialing@outlook.com
// @TermsOfServiceUrl https://github.com/xvjialing
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"school-web/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/article",
			beego.NSInclude(
				&controllers.ArticleController{},
			),
		),

		beego.NSNamespace("/file",
			beego.NSInclude(
				&controllers.FileController{},
			),
		),

		beego.NSNamespace("/token",
			beego.NSInclude(
				&controllers.TokenController{},
			),
		),

		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/teacher",
			beego.NSInclude(
				&controllers.TeacherController{},
			),
		),
		beego.NSNamespace("/leader",
			beego.NSInclude(
				&controllers.LeaderController{},
			),
		),
	)
	beego.AddNamespace(ns)
}

package initialize

import (
	"fmt"

	"fresh-shop/server/global"
	"fresh-shop/server/middleware"
	"fresh-shop/server/plugin/email"
	"fresh-shop/server/utils/plugin"
	"github.com/gin-gonic/gin"
)

func PluginInit(group *gin.RouterGroup, Plugin ...plugin.Plugin) {
	for i := range Plugin {
		PluginGroup := group.Group(Plugin[i].RouterPath())
		Plugin[i].Register(PluginGroup)
	}
}

func InstallPlugin(Router *gin.Engine) {
	PublicGroup := Router.Group("")
	fmt.Println("无鉴权插件安装==》", PublicGroup)
	PrivateGroup := Router.Group("")
	fmt.Println("鉴权插件安装==》", PrivateGroup)
	PrivateGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	//  添加跟角色挂钩权限的插件 示例 本地示例模式于在线仓库模式注意上方的import 可以自行切换 效果相同
	PluginInit(PrivateGroup, email.CreateEmailPlug(
		global.Config.Email.To,
		global.Config.Email.From,
		global.Config.Email.Host,
		global.Config.Email.Secret,
		global.Config.Email.Nickname,
		global.Config.Email.Port,
		global.Config.Email.IsSSL,
	))
}

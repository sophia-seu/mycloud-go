package handler

import (
	"fmt"
	"mycloud/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

//检查是否登录中间件
func CheckLogin(c *gin.Context) {

	//cookie中拿到token
	token, err := c.Cookie("Token")
	if err != nil {
		fmt.Println("cookie", err.Error())
		c.Abort()
	}

	// redis校验token
	openId, err := lib.GetKey(token)
	if err != nil {
		fmt.Println("Get Redis Err:", err.Error())
		c.Redirect(http.StatusFound, "/cloud/login")
		c.Abort()
	}

	// 把openId放入到context内
	c.Set("openId", openId)

	c.Next()

	// user := db.GetUserInfo(openId)

	// if user.Id == 0 {
	// 	//校验失败 返回登录页面
	// 	c.Redirect(http.StatusFound, "/")
	// } else {
	// 	//校验成功 继续执行
	// 	c.Set("openId", openId)
	// 	c.Next()
	// }
}

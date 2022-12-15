package handler

import (
	"fmt"
	"log"
	config "mycloud/conf"
	"mycloud/db"
	dblayer "mycloud/db"
	"mycloud/lib"
	"mycloud/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func Register(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func HandlerRegister(c *gin.Context) {
	username := c.PostForm("username")
	passwd := c.PostForm("password")

	// 检测用户名密码合法性

	if ok := util.CheckUsername(username); ok != true {
		c.JSON(http.StatusOK, gin.H{"msg": "用户名不合法"})
		return
	}

	if ok := util.CheckPassword(passwd); ok != true {
		c.JSON(http.StatusOK, gin.H{"msg": "密码不合法"})
		return
	}

	// 检测用户是否存在
	isExist := db.QueryUserExists(username, passwd)

	if isExist != "not exist" {
		log.Println("用户名已存在!")
		c.JSON(http.StatusOK, gin.H{"msg": "用户名存在"})
		return
	}

	// config := LoadServerConfig()

	// 借用第三方库生成唯一的openId
	openId := xid.New().String()

	dblayer.CreateUser(openId, username, passwd, config.ImageUrl)

	c.Redirect(http.StatusFound, "/cloud/login")
}

func dealRegister(username, passwd string, c *gin.Context) {
	// 检测用户名密码合法性
	if ok := util.CheckUsername(username); ok != true {
		c.JSON(http.StatusOK, gin.H{"msg": "用户名不合法"})
		return
	}

	if ok := util.CheckPassword(passwd); ok != true {
		c.JSON(http.StatusOK, gin.H{"msg": "密码不合法"})
		return
	}
	// 检测用户是否存在
	isExist := db.QueryUserExists(username, passwd)

	if isExist != "not exist" {
		log.Println("用户名已存在!")
		c.JSON(http.StatusOK, gin.H{"msg": "用户名存在"})
		return
	}
	// config := LoadServerConfig()
	// 借用第三方库生成唯一的openId
	openId := xid.New().String()
	dblayer.CreateUser(openId, username, passwd, config.ImageUrl)
	c.Redirect(http.StatusFound, "/cloud/login")
}

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func HandlerLogin(c *gin.Context) {

	// 获得用户名密码
	username := c.PostForm("username")
	passwd := c.PostForm("password")
	mode := c.PostForm("mode")

	if mode == "Sign up" {
		dealRegister(username, passwd, c)
		return
	}

	//校验用户名密码

	res := dblayer.QueryUserExists(username, passwd)
	if res == "not exist" {
		log.Println("用户不存在")
		c.JSON(http.StatusOK, gin.H{"msg": "用户名不存在"})
		return
	} else if res == "not match" {
		log.Println("用户名和密码不匹配")
		c.JSON(http.StatusOK, gin.H{"msg": "用户名密码不匹配"})
		return
	}

	openId := res
	//创建一个token
	hashToken := util.EncodeMd5("token" + string(time.Now().Unix()) + openId)

	//存入redis
	if err := lib.SetKey(hashToken, openId, 24*3600); err != nil {
		fmt.Println("Redis Set Err:", err.Error())
		return
	}

	c.SetCookie("Token", hashToken, config.CookieTime, "/", config.CookieDomain, false, false)

	c.Redirect(http.StatusFound, "/cloud/index")

}

//退出登录
func Logout(c *gin.Context) {
	token, err := c.Cookie("Token")
	if err != nil {
		fmt.Println("cookie", err.Error())
	}

	if err := lib.DelKey(token); err != nil {
		fmt.Println("Del Redis Err:", err.Error())
	}

	c.SetCookie("Token", "", 0, "/", config.CookieDomain, false, false)
	c.Redirect(http.StatusFound, "/cloud/login")
}

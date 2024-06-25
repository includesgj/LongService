package user

import (
	"GinProject12/common"
	sdb "GinProject12/databases"
	"GinProject12/dto"
	"GinProject12/model"
	"GinProject12/response"
	"GinProject12/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func AdminRegister(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	pwd := c.PostForm("password")
	role := c.PostForm("role")

	if len(username) < 3 || len(pwd) < 6 || len(role) == 0 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "数据长度有误请检查")
		return
	}

	// 查询邮箱是否合法
	if ok := util.EmailIsOk(email); !ok {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "邮箱不合法")
		return
	}

	// 查询邮箱是否被注册
	if sdb.FindAdminByEvery("email", email) != nil {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "邮箱已被注册")
		return
	}

	if sdb.FindAdminByEvery("username", username) != nil {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "该用户已被注册")
		return
	}

	// 加密密码
	hashPwd, err := util.HashAndSalt(pwd)

	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "服务器内部问题")
		return
	}

	admin := model.Admin{
		Username: username,
		Password: hashPwd,
		Email:    email,
		Role:     role,
	}

	_, err = sdb.InsertAdminInfo(admin)

	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "服务器内部问题")
		return
	}

	response.Success(c, nil, "注册成功")

}

func AdminLogin(c *gin.Context) {
	username := c.PostForm("username")
	pwd := c.PostForm("password")

	if len(pwd) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码长度不足6位")
		return
	}

	admin := sdb.FindAdminByEvery("username", username)

	if admin == nil {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "该用户不存在")
		return
	}

	ok := util.ComparePasswords(admin.Password, pwd)

	if !ok {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码错误")
		return
	}

	token, err := common.CreateJWT(admin.Username)

	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "服务器内部问题")
		return
	}

	response.Success(c, gin.H{"token": token}, "登陆成功")

}

func AdminInfo(c *gin.Context) {
	user, _ := c.Get("user")
	log.Println(user)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToAdminDto(user.(model.Admin))}})
}

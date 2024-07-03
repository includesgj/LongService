package user

import (
	"GinProject12/common"
	"GinProject12/databases"
	"GinProject12/dto"
	"GinProject12/model"
	"GinProject12/response"
	login_log "GinProject12/serverce/cmd/logs/login-log-panel"
	"GinProject12/util"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// BaseUserRegister 注册
//
//	@Summary		普通用户注册
//	@Description	普通用户注册
//	@Param User body model.User true "request"
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		422
//	@Failure		500
//	@Router			/api/user/register [POST]
func BaseUserRegister(c *gin.Context) {

	req := model.User{}

	if err := util.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	if !util.EmailIsOk(req.Email) {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "邮箱不正确")
		return
	}

	// 查询邮箱是否存在!
	if sdb.FindUserByEvery("email", req.Email) != nil {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "邮箱已存在")
		return
	}

	if len(req.Password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码长度不足6位")
		return
	}

	// 查询用户名是否已存在
	query := sdb.FindUserByEvery("username", req.Username)
	if query != nil {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "该用户已存在")
		return
	}

	hashPassword, err := util.HashAndSalt(req.Password)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "密码加密错误")
		return
	}

	user := &model.User{
		Username: req.Username,
		Password: hashPassword,
		Email:    req.Email,
	}

	sdb.InsertUserInfo(user)
	response.Success(c, nil, "注册成功")

}

// BaseUserLogin 登陆
// @Summary      用户登录
// @Description  用户登录
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        User body model.User true "request"
// @Success      200  {string}  token
// @Failure      422  {string}  err
// @Failure      500  {string}  err
// @Router       /api/user/login [POST]
func BaseUserLogin(c *gin.Context) {

	req := model.User{}

	if err := util.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	user := sdb.FindUserByEvery("username", req.Username)

	if user == nil || user.ID == 0 {
		go login_log.SaveLoginLog(c, errors.New("用户不存在"))
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}

	if !util.ComparePasswords(user.Password, req.Password) {
		go login_log.SaveLoginLog(c, errors.New("密码错误"))
		response.Fail(c, nil, "密码错误")
		return
	}

	// token发送
	token, err := common.CreateJWT(user.Username)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "创建token错误")
		return
	}
	// 返回结果
	go login_log.SaveLoginLog(c, nil)
	response.Success(c, gin.H{"token": token}, "登陆成功")
}

// BaseUserInfo 用户信息
// @Summary      用户信息
// @Description  用户信息
// @Tags         user
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @Success      200
// @Failure      422
// @Failure      500
// @Router       /api/user/info [GET]
func BaseUserInfo(c *gin.Context) {
	user, _ := c.Get("user")
	log.Println(user)

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(model.User))}})
}

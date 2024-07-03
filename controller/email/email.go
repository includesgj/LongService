package email

import (
	"GinProject12/response"
	"GinProject12/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SendEmailCode 向邮箱发送验证码
// @Summary      向邮箱发送验证码
// @Description  向邮箱发送验证码(用query发数据"email")
// @Tags         email
// @Accept       json
// @Param        email query string ture "email"
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /email/code [GET]
func SendEmailCode(c *gin.Context) {
	email := c.Query("email")
	if !util.EmailIsOk(email) {
		response.Response(c, http.StatusBadRequest, 400, nil, "邮箱不正确")
		return
	}

	if err := util.SendEmailCode(email); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}

	response.Success(c, nil, "请在邮箱里查看没有请检查垃圾邮件")

}

// VerifyCode 验证验证码
// @Summary      验证验证码
// @Description  验证验证码(用query发送"email", "code")201代表验证失败
// @Tags         email
// @Accept       json
// @Produce      json
// @Param        email query string ture "email"
// @Param  		 code query string ture "code"
// @Success      200
// @Failure      201
// @Failure      404
// @Failure      500
// @Router       /email/verify [GET]
func VerifyCode(c *gin.Context) {
	email := c.Query("email")
	code := c.Query("code")
	if !util.EmailIsOk(email) {
		response.Response(c, http.StatusBadRequest, 400, nil, "邮箱不正确")
		return
	}

	ok, msg := util.VerificationCode(email, code)

	if !ok {
		response.Response(c, http.StatusOK, 201, nil, msg)
		return
	}
	response.Success(c, nil, msg)

}

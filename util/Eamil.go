package util

import (
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

var (
	userName   = "yzm@songguanju.asia"
	password   = "Sgj123456." // 到时候存到数据库里
	mailServer = "smtp.qiye.aliyun.com:25"
)

const (
	expirationDuration = 5 * time.Minute // 设置验证码的过期时间为5分钟
)

type CodeStore struct {
	Code      string
	Generated time.Time
}

var VCode = make(map[string]*CodeStore)

func sendToMail(to, subject, body string) error {

	// 拼接消息体
	var contentType string

	contentType = "Content-Type: text/plain" + "; charset=UTF-8"

	msg := []byte("To: " + to + "\nFrom: " + userName + "\nSubject: " + subject + "\n" + contentType + "\n\n" + body)

	// 进行身份认证
	hp := strings.Split(mailServer, ":")
	auth := smtp.PlainAuth("", userName, password, hp[0])

	sendTo := strings.Split(to, ";")
	err := smtp.SendMail(mailServer, auth, userName, sendTo, msg)
	return err
}

func SendEmailCode(to string) error {
	code := RandCodeSix()
	body := fmt.Sprintf("亲爱的用户您的验证码是:%s该验证码在5分钟内有效", code)
	if err := sendToMail(to, "龙芯服务验证码", body); err != nil {
		return err
	}

	VCode[to] = &CodeStore{
		Code:      code,
		Generated: time.Now(),
	}

	return nil
}

func VerificationCode(email, code string) (bool, string) {
	info := VCode[email]
	if info == nil {
		return false, "请发送邮件!"
	}
	if info.Code == code && time.Since(info.Generated) > expirationDuration {
		delete(VCode, email)
		return false, "邮件已过期请重新发送!"
	} else if info.Code != code {
		return false, "验证码不正确!"
	}

	delete(VCode, email)
	return true, "成功!"
}

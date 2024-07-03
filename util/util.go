package util

import (
	"GinProject12/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"
)

var (
	VALID = validator.New()
)

func RandomName(n int) string {
	str := []byte("abcdefghijkrmnopqrstuvwxyzABCDEFGHIJKRMNOPQRSTUVWXYZ")
	Name := make([]byte, n)

	for i := 0; i < n; i++ {
		Name[i] = str[rand.Intn(len(str))]
	}

	return string(Name)
}

// HashAndSalt 加密密码
func HashAndSalt(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePasswords 验证密码
func ComparePasswords(hashedPwd string, plainPwd string) bool {
	log.Println(hashedPwd, plainPwd)
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plainPwd))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func EmailIsOk(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	flag := regexp.MustCompile(pattern).MatchString(email)
	return flag
}

// MaxUnit 改数值的最大单位
func MaxUnit(bytes uint64) (float64, string) {
	var unit = []string{"B", "KB", "MB", "GB", "TB"}
	idx := 0
	fByte := float64(bytes)
	for {
		if fByte/1024.0 < 1 || idx >= 5 {
			break
		}
		idx++
		fByte = fByte / 1024
	}
	if int(fByte*1000)%10 > 5 {
		fByte += 0.01
	}
	return fByte, unit[idx]
}

func CheckBindAndValidate(req interface{}, c *gin.Context) error {
	fmt.Println("记得把调用过这个函数的错误处理的response删除")
	if err := c.ShouldBindJSON(req); err != nil {
		response.Fail(c, nil, err.Error())
		return err
	}
	var eStr []string
	if err := VALID.Struct(req); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			eStr = append(eStr, e.Error())
		}
		response.Fail(c, gin.H{"date": eStr}, "失败")
		return err
	}
	return nil
}

func TransformOctal(num int) (os.FileMode, error) {
	mode, err := strconv.ParseUint(fmt.Sprintf("%o", num), 8, 32)
	if err != nil {
		return 0, err
	}
	return os.FileMode(mode), nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandCodeSix() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(100000 + rand.Intn(900000))
}

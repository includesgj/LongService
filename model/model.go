package model

import "time"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email"`
}

type Admin struct {
	ID       int    `json:"id"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type PageInfo struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

// RecycleBin 回收站
type RecycleBin struct {
	Id         int    `json:"id"`
	SourcePath string `json:"sourcePath"`
	DeleteTime string `json:"deleteTime"`
	From       string `json:"from"`
	IsDir      bool   `json:"isDir"`
	Name       string `json:"name"`
	RName      string `json:"rName"`
	Size       int    `json:"size"`
}

// RecoverReq 恢复回收站
type RecoverReq struct {
	Name  string `json:"name"  validate:"required"`
	RName string `json:"rName"  validate:"required"`
	From  string `json:"from"  validate:"required"`
}

type LoginLog struct {
	Id        int    `json:"id"`
	Ip        string `json:"ip"`
	LoginTime string `json:"loginTime"`
	Area      string `json:"area"`
	IsLogin   bool   `json:"isLogin"` // 是否登陆成功
}

type Monitor struct {
	Id          int       `json:"id"`
	CreateUser  string    `json:"createUser"`
	CreateTime  time.Time `json:"createTime"`
	HardWare    string    `json:"hardWare"`
	Detail      string    `json:"detail"`
	Threshold   float64   `json:"threshold"`
	Up          float64   `json:"up"`
	Down        float64   `json:"down"`
	NotifyEmail string    `json:"notifyEmail"`
}

// id createTime 创建时间  TargetDetail 设置巡查内容 createUser 谁创建的
type Patrol struct {
	Id           int       `json:"id"`
	CreateTime   time.Time `json:"createTime"`
	TargetDetail string    `json:"detail"`
	CreateUser   string    `json:"createUser"`
}

type PatrolUser struct {
	PatrolId   int       `json:"patrolId"`
	PatrolTime time.Time `json:"patrolTime"`
	User       string    `json:"patrolUser"`
	Result     bool      `json:"result"`
	Detail     string    `json:"detail"`
}

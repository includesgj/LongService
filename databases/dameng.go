package sdb

import (
	"GinProject12/config"
	"GinProject12/model"
	"database/sql"
	"fmt"
	_ "gitee.com/travelliu/dm"
	_ "github.com/golang/snappy"
	"log"
	"net"
	"strconv"
	"time"
)

/*
 达梦数据库
*/

var (
	DB *sql.DB
)

func GetLocalIPv4() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range interfaces {
		var localAddrs []net.Addr
		localAddrs, err = iface.Addrs()
		if err != nil {
			continue
		}

		for _, add := range localAddrs {
			ip, _, _ := net.ParseCIDR(add.String())
			if ip != nil && ip.To4() != nil && !ip.IsLoopback() {
				return ip.String()
			}
		}
	}

	return ""
}

func InitDM() {
	var err error
	DB, err = sql.Open("dm", fmt.Sprintf("dm://%s:%s@%s:%s", config.DbUser, config.Pwd, config.Addr, config.Port))

	if err != nil {
		panic(err.Error())
	}

	if err = DB.Ping(); err != nil {
		panic(err.Error())
	}
}

func GetDm() *sql.DB {
	var err error
	// addr, err = getLocalIPv4()
	if err != nil {
		panic(err)
	}
	if DB == nil {
		InitDM()
	}
	return DB
}

func CloseDm() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			return
		}
	}
}

func InsertUserInfo(info *model.User) {

	sql := fmt.Sprintf(`INSERT INTO "gorjb"."user" ("username", "password", "email") values('%s', '%s', '%s')`, info.Username, info.Password, info.Email)

	exec, err := DB.Exec(sql)

	if err != nil {
		log.Println(err.Error() + " Insert user info err")
		return
	}

	id, err := exec.LastInsertId()

	if err != nil {
		log.Println(err.Error() + " detail id = " + strconv.FormatInt(id, 10))
		return
	}
	log.Printf("Insert user info successfully id = %d\n", id)
}

func FindUserByEvery(query string, val string) *model.User {
	sql := fmt.Sprintf(`select "id", "username", "password", "email" from "gorjb"."user" where "%s" = '%s'`, query, val)
	row := DB.QueryRow(sql)
	if row.Err() != nil {
		log.Println(row.Err().Error() + " find user by id fail")
		return nil
	}

	var user model.User

	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		log.Println(err.Error() + " scan fail")
		return nil
	}
	return &user
}

func FindAdminByEvery(query string, val string) *model.Admin {
	sql := fmt.Sprintf(`SELECT "id", "username", "password", "email", "role" FROM "gorjb"."admin" WHERE "%s" = '%s'`, query, val)
	exec := DB.QueryRow(sql)

	if exec.Err() != nil {
		log.Println(exec.Err().Error() + " find admin by id fail")
		return nil
	}

	var admin model.Admin
	err := exec.Scan(&admin.ID, &admin.Username, &admin.Password, &admin.Email, &admin.Role)

	if err != nil {
		log.Println(err.Error() + " scan fail")
		return nil
	}
	return &admin
}

func InsertAdminInfo(info model.Admin) (int64, error) {
	sql := fmt.Sprintf(`INSERT INTO "gorjb"."admin" ("username", "password","email" , "role") values('%s', '%s', '%s', '%s')`, info.Username, info.Password, info.Email, info.Role)
	exec, err := DB.Exec(sql)

	if err != nil {
		log.Println(err.Error() + "add admin fail")
		return -1, err
	}

	id, err := exec.LastInsertId()

	if err != nil {
		log.Println(err.Error())
		return -1, err
	}
	log.Printf("Insert admin successfully id = %d\n", id)

	return id, nil

}

func InsertRecycleBinInfo(info model.RecycleBin) (int64, error) {

	isDir := 0
	if info.IsDir {
		isDir = 1
	}

	sql := fmt.Sprintf(`INSERT INTO "gorjb"."RecycleBin" ("name", "rName", "sourcePath", "from", "isDir", "deleteTime", "size") values('%s', '%s', '%s', '%s', '%v', '%v', '%d')`, info.Name, info.RName, info.SourcePath, info.From, isDir, info.DeleteTime, info.Size)
	exec, err := DB.Exec(sql)
	if err != nil {
		log.Println(err.Error())
		return -1, err
	}
	id, err := exec.LastInsertId()

	if err != nil {
		log.Println(err.Error())
		return -1, err
	}

	return id, nil
}

func RecycleBinPage(page model.PageInfo) ([]model.RecycleBin, error) { // ORDER BY your_column OFFSET 10 ROWS FETCH NEXT 10 ROWS ONLY;
	sql := fmt.Sprintf(`SELECT "id", "deleteTime", "name", "rName", "sourcePath", "from", "size", "isDir" FROM "gorjb"."RecycleBin" ORDER BY 'deleteTime' OFFSET %d ROWS FETCH NEXT %d ROWS ONLY`, page.Page-1, page.PageSize)

	exec, err := DB.Query(sql)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if exec.Err() != nil {
		log.Println(err)
		return nil, err
	}

	var list []model.RecycleBin

	for exec.Next() {
		var info model.RecycleBin
		err = exec.Scan(&info.Id, &info.DeleteTime, &info.Name, &info.RName, &info.SourcePath, &info.From, &info.Size, &info.IsDir)

		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		list = append(list, info)

	}

	return list, nil
}

// RecycleBinInfo 查询回收站的某一条
func RecycleBinInfo(req model.RecoverReq) (*model.RecycleBin, error) {
	sql := fmt.Sprintf(`SELECT "id", "deleteTime", "name", "rName", "sourcePath", "from", "size", "isDir" FROM "gorjb"."RecycleBin" WHERE 'from' = '%s' AND 'name' = '%s' AND 'rName' = '%s'`, req.From, req.Name, req.Name)
	row := DB.QueryRow(sql)

	if row.Err() != nil {
		return nil, row.Err()
	}

	var info model.RecycleBin
	if err := row.Scan(&info.Id, &info.DeleteTime, &info.Name, &info.RName, &info.SourcePath, &info.From, &info.Size, &info.IsDir); err != nil {
		return nil, err
	}

	return &info, nil

}

func DelRecycleBin(id int) error {
	sql := fmt.Sprintf(`DELETE FROM "gorjb"."RecycleBin" WHERE "id" = '%d'`, id)
	_, err := DB.Exec(sql)
	if err != nil {
		return err
	}
	log.Println("删除成功")
	return nil
}

func InsertLoginLog(info model.LoginLog) error {
	var is = 0
	if info.IsLogin {
		is = 1
	}
	sql := fmt.Sprintf(`INSERT INTO "gorjb"."LoginLog" ("ip", "area", "loginTime", "isLogin") values('%s', '%s', '%s', '%d')`, info.Ip, info.Area, info.LoginTime, is)
	exec, err := DB.Exec(sql)
	if err != nil {
		return err
	}
	id, err := exec.LastInsertId()

	if err != nil {
		return err
	}

	log.Printf("插入成功id=%d\n", id)
	return nil
}

func LoginLogPage(page model.PageInfo) ([]model.LoginLog, error) {
	sql := fmt.Sprintf(`SELECT "id", "ip", "area", "loginTime", "isLogin" from "gorjb"."LoginLog" ORDER BY 'loginTime' OFFSET %d ROWS FETCH NEXT %d ROWS ONLY`, page.Page-1, page.PageSize)
	exec, err := DB.Query(sql)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if exec.Err() != nil {
		log.Println(err)
		return nil, err
	}

	var list []model.LoginLog

	for exec.Next() {
		var info model.LoginLog
		err = exec.Scan(&info.Id, &info.Ip, &info.Area, &info.LoginTime, &info.IsLogin)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		list = append(list, info)

	}

	return list, nil
}

func InsertMonitorInfo(info model.Monitor) (int, error) {

	sql := fmt.Sprintf(`INSERT INTO "gorjb"."Monitor" ("createTime", "createUser", "hardWare", "threshold", "detail", "up", "down", "notifyEmail") values('%s', '%s', '%s', '%.2f', '%s', '%.2f', '%.2f', '%s')`, info.CreateTime.String(), info.CreateUser, info.HardWare, info.Threshold, info.Detail, info.Up, info.Down, info.NotifyEmail)

	exec, err := DB.Exec(sql)
	if err != nil {
		return -1, err
	}
	id, err := exec.LastInsertId()

	if err != nil {
		return -1, err
	}

	log.Printf("插入成功id=%d\n", id)
	return int(id), nil
}

func DelMonitorInfo(id int) error {
	sql := fmt.Sprintf(`DELETE FROM "gorjb"."Monitor" WHERE "id" = '%d'`, id)
	_, err := DB.Exec(sql)
	if err != nil {
		return err
	}
	log.Println("删除成功")
	return nil
}

func SelectMonitor() ([]model.Monitor, error) {
	sql := fmt.Sprintf(`SELECT "id", "createTime", "createUser", "hardWare", "threshold", "detail", "up", "down", "notifyEmail" FROM "gorjb"."Monitor"`)

	exec, err := DB.Query(sql)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if exec.Err() != nil {
		log.Println(err)
		return nil, err
	}

	var list []model.Monitor

	for exec.Next() {
		var info model.Monitor
		var strTime string
		err = exec.Scan(&info.Id, &strTime, &info.CreateUser, &info.HardWare, &info.Threshold, &info.Detail, &info.Up, &info.Down, &info.NotifyEmail)
		info.CreateTime, _ = time.Parse("2006-01-02 15:04:05", strTime)

		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		list = append(list, info)
	}
	return list, nil
}

func SelectMonitorPage(page model.PageInfo) ([]model.Monitor, error) {
	sql := fmt.Sprintf(`SELECT "id", "createTime", "createUser", "hardWare", "threshold", "detail", "up", "down", "notifyEmail" FROM "gorjb"."Monitor" ORDER BY 'createTime' OFFSET %d ROWS FETCH NEXT %d ROWS ONLY`, page.Page-1, page.PageSize)

	exec, err := DB.Query(sql)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if exec.Err() != nil {
		log.Println(err)
		return nil, err
	}

	var list []model.Monitor

	for exec.Next() {
		var info model.Monitor
		var strTime string
		err = exec.Scan(&info.Id, &strTime, &info.CreateUser, &info.HardWare, &info.Threshold, &info.Detail, &info.Up, &info.Down, &info.NotifyEmail)
		info.CreateTime, _ = time.Parse("2006-01-02 15:04:05", strTime)

		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		list = append(list, info)
	}
	return list, nil
}

/*
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
*/
func InsertPatrol(info model.Patrol) (int, error) {
	sql := fmt.Sprintf(`INSERT INTO "gorjb"."Patrol" ("createTime", "targetDetail", "createUser") values('%s', '%s', '%s')`, info.CreateTime, info.TargetDetail, info.CreateUser)

	exec, err := DB.Exec(sql)
	if err != nil {
		return -1, err
	}
	id, err := exec.LastInsertId()

	if err != nil {
		return -1, err
	}

	log.Printf("插入成功id=%d\n", id)
	return int(id), nil
}

func SelectPatrol(page model.PageInfo) ([]model.Patrol, error) {
	sql := fmt.Sprintf(`SELECT "id", "createTime", "targetDetail", "createUser" FROM "gorjb"."Patrol"  ORDER BY 'createTime' OFFSET %d ROWS FETCH NEXT %d ROWS ONLY`, page.Page-1, page.PageSize)

	exec, err := DB.Query(sql)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if exec.Err() != nil {
		log.Println(err)
		return nil, err
	}

	var list []model.Patrol

	for exec.Next() {
		var info model.Patrol
		var strTime string
		err = exec.Scan(&info.Id, &strTime, &info.TargetDetail, &info.CreateUser)
		info.CreateTime, _ = time.Parse("2006-01-02 15:04:05", strTime)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		list = append(list, info)
	}
	return list, nil
}

func DelPatrol(id int) error {
	sql := fmt.Sprintf(`DELETE FROM "gorjb"."Patrol" WHERE "id" = '%d'`, id)
	_, err := DB.Exec(sql)
	if err != nil {
		return err
	}
	log.Println("删除成功")
	return nil
}

/*
	type PatrolUser struct {
		PatrolId   int       `json:"patrolId"`
		PatrolTime time.Time `json:"patrolTime"`
		User       string    `json:"patrolUser"`
		Result     bool      `json:"result"`
		Detail     string    `json:"detail"`
	}
*/
func InsertPatrolUser(info model.PatrolUser) (int, error) {
	var is = 0
	if info.Result {
		is = 1
	}

	sql := fmt.Sprintf(`INSERT INTO "gorjb"."PatrolUser" ("patrolId", "patrolTime", "user", "result", "detail") values('%d', '%s', '%s', '%d', '%s')`, info.PatrolId, info.PatrolTime, info.User, is, info.Detail)

	exec, err := DB.Exec(sql)
	if err != nil {
		return -1, err
	}
	id, err := exec.LastInsertId()

	if err != nil {
		return -1, err
	}

	log.Printf("插入成功id=%d\n", id)
	return int(id), nil
}

func SelectPatrolUser(page model.PageInfo) ([]model.PatrolUser, error) {
	sql := fmt.Sprintf(`SELECT "patrolId", "patrolTime", "user", "result", "detail" FROM "gorjb"."PatrolUser"  ORDER BY 'patrolTime' OFFSET %d ROWS FETCH NEXT %d ROWS ONLY`, page.Page-1, page.PageSize)

	exec, err := DB.Query(sql)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if exec.Err() != nil {
		log.Println(err)
		return nil, err
	}

	var list []model.PatrolUser

	for exec.Next() {
		var info model.PatrolUser
		err = exec.Scan(&info.PatrolId, &info.PatrolTime, &info.User, &info.Result, &info.Detail)

		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		list = append(list, info)
	}
	return list, nil
}

func DelPatrolUser(id int) error {
	sql := fmt.Sprintf(`DELETE FROM "gorjb"."PatrolUser" WHERE "id" = '%d'`, id)
	_, err := DB.Exec(sql)
	if err != nil {
		return err
	}
	log.Println("删除成功")
	return nil
}

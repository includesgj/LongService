package sdb

import (
	"GinProject12/model"
	"database/sql"
	"errors"
	"fmt"
	_ "gitee.com/travelliu/dm"
	_ "github.com/golang/snappy"
	"log"
	"net"
	"strconv"
)

/*
 达梦数据库
*/

var (
	DB     *sql.DB
	dbUser = "SYSDBA"    // dmdba
	pwd    = "SYSDBA001" // 123456
	addr   = "127.0.0.1" // "192.168.10.105" //"192.168.1.150" // "172.16.102.211" // 学校109
	port   = "5236"
)

func getLocalIPv4() (string, error) {
	fmt.Println("如果数据库不在本机上就把ip改到运行机上把这条删除!")
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
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
				return ip.String(), nil
			}
		}
	}

	return "", errors.New("no IPv4 address found")
}

func InitDM() {
	var err error
	DB, err = sql.Open("dm", fmt.Sprintf("dm://%s:%s@%s:%s", dbUser, pwd, addr, port))

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

//// AllStudents 给老师看的 查询所有学生
//func AllStudents() []model.Student {
//	query, err := DB.Query(`select "id", "name", "student_id", "password", "class", "telephone" from "gorjb"."student"`)
//
//	if err != nil {
//		log.Println(err.Error() + " select student error")
//		return nil
//	}
//
//	defer query.Close()
//
//	var student []model.Student
//
//	for query.Next() {
//		var stu model.Student
//		err = query.Scan(&stu.ID, &stu.Name, &stu.StudentId, &stu.Class, &stu.Password, &stu.Telephone)
//
//		if err != nil {
//			log.Println(err.Error())
//			return nil
//		}
//		student = append(student, stu)
//	}
//	return student
//
//}

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

// DeleteUserById 给老师用的 删除用户
func DeleteUserById(id int) {
	sql := fmt.Sprintf(`DELETE FROM "gorjb"."user" WHERE "id" = '%d'`, id)
	_, err := DB.Exec(sql)

	if err != nil {
		log.Println(err.Error() + " delete user fail")
		return
	}
	log.Printf("Delete user id is %d\n", id)
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
	sql := fmt.Sprintf(`SELECT "id", "ip", "area", "loginTime", "isLogin" form "gorjb"."LoginLog" ORDER BY 'deleteTime' OFFSET %d ROWS FETCH NEXT %d ROWS ONLY`, page.Page-1, page.PageSize)
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

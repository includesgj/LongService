// Copyright © 2021 Bin Liu <bin.liu@enmotech.com>

package dm

import (
	"database/sql"
	"fmt"
	"testing"
)

func Test_Connect(t *testing.T) {
	dsn := "dm://sysdba:SYSDBA@172.23.1.54:52360?appName=mtk&connectTimeout=3000&logLevel=all"
	db, err := sql.Open("dm", dsn)
	if err != nil {
		t.Error(err)
		return
	}
	sqlText := "select sf_get_unicode_flag()"
	var charSet string

	err = db.QueryRow(sqlText).Scan(&charSet)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(charSet)
}

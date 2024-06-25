package GoTest

import (
	sdb "GinProject12/databases"
	"GinProject12/model"
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {
	dm := sdb.GetDm()
	defer dm.Close()

	info := model.PageInfo{Page: 1, PageSize: 1}

	page, err := sdb.RecycleBinPage(info)

	if err != nil {
		fmt.Println(err)
	}

	for i, j := range page {
		fmt.Println(i, j)
	}

}

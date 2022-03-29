package week02

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"testing"
)

func Test1(t *testing.T) {
	week02()
}

//week02 sql.ErrNoRows的错误处理
func week02() {
	if err := doStore(); err != nil {
		fmt.Printf("%+v", err)
	}
}

//doStore 执行store层逻辑
func doStore() error {
	err := sql.ErrNoRows
	query := "select name from users where id = ?"
	return errors.Wrap(err, query)
}

package DBOpt

import (
	"database/sql"
	"errors"

	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

type BaseDB struct {
	//mMysqlConn *sql.DB
	//mTimeOut   time.Duration
	mAddr string

	errDBConnect    error
	errOptException error
}

func (opt *BaseDB) InitDatabase(addr string) error {
	opt.mAddr = addr

	opt.errDBConnect = errors.New("数据连接失败")
	opt.errOptException = errors.New("数据库异常，操作失败")
	return nil
}

func (opt *BaseDB) connectDB() (*sql.DB, error) {
	conn, err := sql.Open("mysql", opt.mAddr)
	return conn, err
}
func (opt *BaseDB) releaseDB(conn *sql.DB) {
	if conn != nil {
		conn.Close()
	}
}
func (opt *BaseDB) exec(conn *sql.DB, sqlString string, args ...interface{}) (err error) {
	if len(opt.mAddr) == 0 {
		//log.Error("err~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~没有配置数据库地址")
		return nil
	}
	if conn == nil {
		conn, err = opt.connectDB()
		defer opt.releaseDB(conn)
		if err != nil {
			return err
		}
	}
	stmt, err := conn.Prepare(sqlString)
	if err != nil {
		log.Error("err:", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(args...)
	if err != nil {
		log.Error("err:", err)
		return err
	}

	return err
}

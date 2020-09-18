package sql_factory

import (
	"database/sql"
	"fmt"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/yml_config"
	"goskeleton/app/core/event_manage"
	"time"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

var mysqlDriverWrite *sql.DB
var mysqlDriverRead *sql.DB

var sqlServerDriverWrite *sql.DB
var sqlServerDriverRead *sql.DB

var postgreDriverWrite *sql.DB
var postgreDriverRead *sql.DB

func initSqlServer(sqlType, readOrWrite string) *sql.DB {
	configFac := yml_config.CreateYamlFactory()
	var tmpSqlType string
	var tmpDriver *sql.DB
	var err error
	switch sqlType {
	case "mysql":
		tmpSqlType = "Mysql"
	case "sqlserver", "mssql":
		tmpSqlType = "SqlServer"
	case "postgre", "postgres", "postgresql":
		tmpSqlType = "PostgreSql"
	default:
		return nil
	}
	// 初始化相同配置
	Host := configFac.GetString(tmpSqlType + "." + readOrWrite + ".Host")
	Port := configFac.GetString(tmpSqlType + "." + readOrWrite + ".Port")
	User := configFac.GetString(tmpSqlType + "." + readOrWrite + ".User")
	Pass := configFac.GetString(tmpSqlType + "." + readOrWrite + ".Pass")
	DataBase := configFac.GetString(tmpSqlType + "." + readOrWrite + ".DataBase")
	SetMaxIdleConns := configFac.GetInt(tmpSqlType + "." + readOrWrite + ".SetMaxIdleConns")
	SetMaxOpenConns := configFac.GetInt(tmpSqlType + "." + readOrWrite + ".SetMaxOpenConns")
	SetConnMaxLifeTime := configFac.GetDuration(tmpSqlType + "." + readOrWrite + "SetConnMaxLifetime")

	if sqlType == "mysql" {
		Charset := configFac.GetString(tmpSqlType + "." + readOrWrite + ".Charset")
		SqlConnString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?ParseTime=True&loc=Local&charset=%s", User, Pass, Host, Port, DataBase, Charset)
		switch readOrWrite {
		case "Write", "Read":
			tmpDriver, err = sql.Open("mysql", SqlConnString)
		default:
			variable.ZapLog.Error(my_errors.ErrorsDbSqlWriteReadInitFail + readOrWrite)
			return nil
		}
		if err != nil {
			variable.ZapLog.Error(my_errors.ErrorsDbSqlDriverInitFail, zap.Error(err))
			return nil
		}
		tmpDriver.SetMaxIdleConns(SetMaxIdleConns)
		tmpDriver.SetMaxOpenConns(SetMaxOpenConns)
		tmpDriver.SetConnMaxLifetime(SetConnMaxLifeTime * time.Second)
		// 将需要销毁的时间统一注册在事件管理器，由程序退出时统一注销
		event_manage.CreateEventManageFactory().Set(
			variable.EventDestroyPrefix+tmpSqlType+readOrWrite,
			func(args ...interface{}) {
				_ = tmpDriver.Close()
			},
		)
		switch readOrWrite {
		case "Write":
			mysqlDriverWrite = tmpDriver
		case "Read":
			mysqlDriverRead = tmpDriver
		default:
			return nil
		}
		return tmpDriver
	}
	return nil
}

func GetOneSqlClient(sqlType, readOrWrite string) *sql.DB {
	if !strings.Contains("mysql,sqlserver,mssql,postgre,postgres,postgresql", sqlType) {


		variable.ZapLog.Error(my_errors.ErrorsDbDriverNotExists + sqlType)
		return nil
	}
	if !strings.Contains("Read,Write", readOrWrite) {
		variable.ZapLog.Error(my_errors.ErrorsDbSqlWriteReadInitFail + "," + readOrWrite)
		return nil
	}
	var maxRetryTimes int
	var reConnectInterval time.Duration
	configFac := yml_config.CreateYamlFactory()

	var dbDriver *sql.DB
	switch sqlType {
	case "mysql":
		if readOrWrite == "Write" {
			if mysqlDriverWrite == nil {
				dbDriver = initSqlServer(sqlType, readOrWrite)
			} else {
				dbDriver = mysqlDriverWrite
			}
		} else if readOrWrite == "Read" {
			if mysqlDriverRead == nil {
				dbDriver = initSqlServer(sqlType, readOrWrite)
			} else {
				dbDriver = mysqlDriverRead
			}
		}
		maxRetryTimes = configFac.GetInt("Mysql." + readOrWrite + ".PingFailRetryTimes")
		reConnectInterval = configFac.GetDuration("Mysql." + readOrWrite + ".ReConnectInterval")
	default:
		variable.ZapLog.Error(my_errors.ErrorsDbDriverNotExists + "," + sqlType)
		return nil
	}

	for i := 1; i <= maxRetryTimes; i ++ {
		// ping 失败允许重试
		if err := dbDriver.Ping(); err != nil { // 获取一个连接失败，进行重试
			dbDriver = initSqlServer(sqlType, readOrWrite)
			time.Sleep(time.Second * reConnectInterval)
			if i == maxRetryTimes {
				variable.ZapLog.Error(sqlType + my_errors.ErrorsDbGetConnFail, zap.Error(err))
				return nil
			}
		} else {
			break
		}
	}
	return dbDriver
}

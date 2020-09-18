package model

import (
	"database/sql"
	"strings"
	"goskeleton/app/utils/sql_factory"
	"goskeleton/app/global/variable"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/utils/yml_config"
)

// 
func CreateBaseSqlFactory(sqlType string) *BaseModel {
	var dbType string
	var sqlDriverRead *sql.DB
	
	sqlType = strings.ToLower(strings.Replace(sqlType, " ", "", -1))
	sqlDriverWrite := sql_factory.GetOneSqlClient(sqlType, "Write")

	switch sqlType {
	case "mysql":
		dbType = "Mysql"
	case "sqlserver", "mssql":
		dbType = "SqlServer"
	case "postgre", "postgres", "postgresql":
		dbType = "PostgreSql"
	default:
		variable.ZapLog.Error(my_errors.ErrorsDbDriverNotExists + sqlType)
		return nil
	}
	// 配置项是否开启读写分离
	isOpenReadDb := yml_config.CreateYamlFactory().GetInt(dbType + ".IsOpenReadDb")
	if isOpenReadDb == 1 {
		sqlDriverRead = sql_factory.GetOneSqlClient(sqlType, "Read")
	} else {
		sqlDriverRead = sqlDriverWrite
	}
	return &BaseModel{
		dbDriverWrite: sqlDriverWrite,
		dbDriverRead: sqlDriverRead,
	}
}

//
type BaseModel struct {
	dbDriverWrite *sql.DB
	dbDriverRead  *sql.DB
	stm           *sql.Stmt
}


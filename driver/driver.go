package driver

import (
	"Auth/entity"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)


func ConnectToDB(config entity.MySQLConfig) (*sql.DB,error){

	db,err := sql.Open("mysql",config.DbUser+":"+config.DbPass+"@/"+config.DbName)

	if err != nil{
		return nil,err
	}
	return db,nil
}

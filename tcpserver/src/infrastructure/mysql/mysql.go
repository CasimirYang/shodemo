package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"tcpserver/infrastructure/po"
	"tcpserver/trace"
	"time"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/spo_db")
	if err != nil {
		_ = trace.Logger.Error(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(30)
	db.SetMaxIdleConns(10)
}

func GetUser(userName, password string) (*po.UserPO, error) {
	rows := db.QueryRow("select id,user_name,nick_name,password,profile from user_base_info_tab where user_name = ? and password= ? ", userName, password)
	var user po.UserPO
	err := rows.Scan(&user.Id, &user.UserName, &user.NickName, &user.Password, &user.Profile)
	return &user, err
}

func GetUserByUid(uid int64) (*po.UserPO, error) {
	rows := db.QueryRow("select id,user_name,nick_name,password,profile from user_base_info_tab where id= ? ", uid)
	var user po.UserPO
	err := rows.Scan(&user.Id, &user.UserName, &user.NickName, &user.Password, &user.Profile)
	return &user, err
}

func UpdateNickName(uid int64, nickName string) error {
	_, err := db.Exec("update user_base_info_tab set nick_name=? where id=? ", nickName, uid)
	return err
}

func UpdateProfile(uid int64, profile string) error {
	_, err := db.Exec("update user_base_info_tab set profile=? where id=? ", profile, uid)
	return err
}

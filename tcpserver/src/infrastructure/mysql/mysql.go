package mysql

import (
	"database/sql"
	commonLog "git.garena.com/jinghua.yang/entry-task-common/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"os"
	"tcpserver/infrastructure/po"
	"time"
)

var db *sql.DB

func init() {
	var err error
	jdbc := viper.GetString("mysql.jdbc")
	db, err = sql.Open("mysql", jdbc)
	if err != nil {
		commonLog.SugarLogger.Error(err)
		os.Exit(1)
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

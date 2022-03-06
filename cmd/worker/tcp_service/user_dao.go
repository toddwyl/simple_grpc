package tcp_service

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"simple_grpc/proto"
)

type TUser struct {
	Id         int64  `db:"id"`
	Username   string `db:"username"`
	Nickname   string `db:"nickname"`
	Password   string `db:"password"`
	ProfilePic string `db:"profile_pic"`
}

// insert t_user
var sqlInsertUser = `INSERT INTO t_user(username, nickname, password, profile_pic, ctime, mtime) 
	VALUES(?, ?, CONCAT('*', UPPER(SHA1(UNHEX(SHA1(?))))), ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

// update profile
var sqlUpdateProfile = ` UPDATE t_user SET profile_pic=?, mtime=CURRENT_TIMESTAMP WHERE username=? `

// update nickname
var sqlUpdateNickname = ` UPDATE t_user SET nickname=?, mtime=CURRENT_TIMESTAMP WHERE username=? `

// select by username
var sqlSelectByUsername = ` SELECT id,username,nickname,password,profile_pic FROM t_user WHERE username=? LIMIT 1 `

// 密码校验
var sqlSelectPassword = ` SELECT CONCAT('*', UPPER(SHA1(UNHEX(SHA1(?))))) `

// for test
var sqlSelectUsernameList = ` SELECT username FROM t_user LIMIT ? `

// selectUsernameList 查询出size个用户名
func selectUsernameList(db *sqlx.DB, size int) ([]string, error) {
	rows, err := db.Queryx(sqlSelectUsernameList, size)
	if err != nil {
		return make([]string, 0, 0), err
	}
	defer rows.Close()
	var results []string
	for rows.Next() {
		var item string
		err := rows.Scan(&item)
		if err != nil {
			logrus.Errorf("scan username row err:%v", err)
			return results, err
		}
		results = append(results, item)
	}
	return results, nil
}

// insertUser 插入用户记录
// return 返回状态码，1-成功，0-失败
func insertUser(db *sqlx.DB, request *proto.RegisterRequest) (int, error) {
	result, err := db.Exec(sqlInsertUser, request.Username, request.Nickname, request.Password, request.ProfilePic)
	fmt.Println(request)
	fmt.Println(request.Password)
	if err != nil {
		// eg: Duplicate entry 'xxx' for key 't_user.unique_idx_username'
		logrus.Warnf("insert user err:%v, request:%s", err, request)
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		logrus.Errorf("insert user err: %v, sql:%s", err, sqlInsertUser)
		return 0, err
	}
	lastId, errId := result.LastInsertId()
	if errId != nil {
		logrus.Errorf("insert user err: %v, sql:%s", err, sqlInsertUser)
		return 0, errId
	}
	logrus.Infof("insert result, rows:%v, last id:%v \n", rows, lastId)
	return 1, nil
}

// selectPassword mysql密码加密
func selectPassword(db *sqlx.DB, password string) string {
	rows, err := db.Queryx(sqlSelectPassword, password)
	if err != nil {
		logrus.Errorf("select password err:%v", err)
		return ""
	}
	defer rows.Close()
	if rows.Next() {
		var encoded string
		err := rows.Scan(&encoded)
		if err == nil {
			return encoded
		}
	}
	logrus.Errorf("scan result err:%v ", err)
	return ""
}

// queryUserByUsername 根据username查询用户信息
func queryUserByUsername(db *sqlx.DB, username string) (*TUser, error) {
	rows, err := db.Queryx(sqlSelectByUsername, username)
	if err != nil {
		logrus.Errorf("query user by username:%s err:%v", username, err)
		return nil, fmt.Errorf("username:%s not exists", username)
	}
	defer rows.Close()
	// SQL限定只返回一行
	var user TUser
	if rows.Next() {
		errMapper := rows.StructScan(&user)
		if errMapper != nil {
			return nil, fmt.Errorf("query user:%s result mapper err:%v", username, errMapper)
		}
		logrus.Debugf("query user:%s by username, get:%v ", username, user)
		return &user, nil
	} else {
		return nil, fmt.Errorf("query user:%s not exists", username)
	}
}

//
// updateUserProfile 更新用户头像
func updateUserProfile(db *sqlx.DB, ProfilePic string, Username string) (int, error) {
	res, err := db.Exec(sqlUpdateProfile, ProfilePic, Username)
	if err != nil {
		logrus.Errorf("update user:%s profile failure:%v", Username, err)
		return 0, err
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		logrus.Errorf("update user:%s profile failure, rows:%v, err:%v", Username, rows, err)
		return 0, err
	}
	logrus.Infof("update user:%s profile success", Username)
	return 1, nil
}

// updateUserNick 更新用户昵称
func updateUserNick(db *sqlx.DB, Nickname string, Username string) (int, error) {
	res, err := db.Exec(sqlUpdateNickname, Nickname, Username)
	if err != nil {
		logrus.Errorf("update user:%s nickname failure:%v", Username, err)
		return 0, err
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		logrus.Errorf("update user:%s nickname failure, rows:%v, err:%v", Username, rows, err)
		return 0, err
	}
	logrus.Infof("update user:%s nickname success", Username)
	return 1, nil
}

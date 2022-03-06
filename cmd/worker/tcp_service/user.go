package tcp_service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"simple_grpc/cmd/global"
	"simple_grpc/proto"
	"strings"
	"time"
)

type UserService struct {
	ctx   context.Context
	db    *sqlx.DB
	cache *redis.Client
	proto.UnimplementedUserServiceServer
}

func NewUserService(ctx context.Context) UserService {
	svc := UserService{ctx: ctx, db: global.DBEngine, cache: global.CacheClient}
	return svc
}

// Login 用户登录
// return user facade.User, err int
// user: 	成功则返回用户信息
// token:	登录成功之后返回token
func (svc *UserService) Login(ctx context.Context, request *proto.LoginRequest) (*proto.LoginReply, error) {
	// 校验账号密码是否正确
	// if not exist, query from db, and refresh into redis
	userT, daoErr := queryUserByUsername(svc.db, request.Username)
	if daoErr != nil {
		logrus.Warnf("query user from db by username:%s failure, msg:%v \n", request.Username, daoErr)
		return nil, fmt.Errorf("svc.UserLogin: username not exist")
	}
	encoded := selectPassword(svc.db, request.Password)
	fmt.Println("userT:", userT)
	fmt.Println("request:", request)
	fmt.Println("encoded:", encoded)
	if strings.Compare(encoded, userT.Password) != 0 { // 密码不正确
		logrus.Warnf("login failure, password mismatch username:%s", request.Username)
		return nil, fmt.Errorf("svc.UserLogin: pwd incorrect")
	}
	//sessionID := generateToken(request.Username)
	sessionID := uuid.NewV4().String()
	//go func() {
	svc.setCacheTokenForUser(ctx, sessionID, userT.Username)
	svc.cacheUserInfoIntoRedis(ctx, userT)
	//}()
	return &proto.LoginReply{SessionId: sessionID}, nil
}

// Register 用户注册
// return username string, err facade.BizError
// username:	注册成功后的用户名（如果传入的用户名带有前或后空格，则会去除空格后再注册并返回）
func (svc *UserService) Register(ctx context.Context, request *proto.RegisterRequest) (*proto.RegisterReply, error) {
	if request.Username == "" || request.Nickname == "" || request.Password == "" {
		// todo: 校验字段长度和字符规范
		return nil, fmt.Errorf("svc.Register Username or Nickname or Password should not be empty")
	}
	// insert user
	_, errInsert := insertUser(svc.db, request)
	if errInsert != nil {
		// failure
		return nil, fmt.Errorf("svc.Register insert error")
	}
	go func() {
		user, err := queryUserByUsername(svc.db, request.Username)
		if err == nil {
			svc.cacheUserInfoIntoRedis(ctx, user)
		}
	}()
	return &proto.RegisterReply{}, nil
}
func (svc UserService) GetUser(ctx context.Context, request *proto.GetUserRequest) (*proto.GetUserReply, error) {
	// Get Username
	logrus.Debugf("request.SessionId:%v", request.SessionId)
	username, err := svc.queryCacheTokenForUser(ctx, request.SessionId)
	if err != nil {
		logrus.Errorf("svc.queryCacheTokenForUser failed: %v", err)
		return nil, err
	}
	// Try loading user profile from cache
	cacheProfileKey := global.UserProfileCachePrefix + username
	userProfCache, errCache := svc.cache.Get(ctx, cacheProfileKey).Result()
	if errCache == nil {
		userGetCacheResp := &proto.GetUserReply{}
		err = json.Unmarshal([]byte(userProfCache), userGetCacheResp)
		if err != nil {
			logrus.Errorf("svc.UserGet: Unmarshal cache failed: %v", err)
		} else {
			return userGetCacheResp, nil
		}
	}

	// Query user from DB
	user, errQuery := queryUserByUsername(svc.db, username)
	if errQuery != nil {
		logrus.Errorf("queryUserByUsername failed: %v", errQuery)
		return nil, errQuery
	}
	go svc.cacheUserInfoIntoRedis(ctx, user)
	return TUser2prototype(user), nil
}
func (svc UserService) EditUser(ctx context.Context, request *proto.EditUserRequest) (*proto.EditUserReply, error) {
	// update db
	username, err := svc.queryCacheTokenForUser(ctx, request.SessionId)
	if err != nil {
		return nil, err
	}
	_, errNick := updateUserNick(svc.db, request.Nickname, username)
	if errNick != nil {
		logrus.Errorf("update user:%s Nickname failure:%v", username, err)
		return nil, errNick
	}
	_, errProfile := updateUserProfile(svc.db, request.ProfilePic, username)
	if errProfile != nil {
		logrus.Errorf("update user:%s ProfilePic failure:%v", username, err)
		return nil, errProfile
	}

	user, err := queryUserByUsername(svc.db, username)
	if err == nil {
		svc.cacheUserInfoIntoRedis(ctx, user)
	}
	userReply := &proto.EditUserReply{}
	return userReply, nil

}

func TUser2prototype(user *TUser) *proto.GetUserReply {
	userGetResp := &proto.GetUserReply{
		Username:   user.Username,
		Nickname:   user.Nickname,
		ProfilePic: user.ProfilePic,
	}
	return userGetResp
}

// cacheUserInfoIntoRedis 缓存用户信息到redis
// 缓存时间： 30分钟
func (svc UserService) cacheUserInfoIntoRedis(ctx context.Context, user *TUser) (*proto.GetUserReply, error) {
	// set into redis
	cacheKey := global.UserProfileCachePrefix + user.Username
	userGetResp := TUser2prototype(user)
	// Set user to cache
	cacheUser, errJson := json.Marshal(userGetResp)
	if errJson != nil {
		logrus.Errorf("errJson: marshal cache failed: %v", errJson)
	}
	err := svc.cache.Set(ctx, cacheKey, cacheUser, 30*time.Minute).Err() // Omit error
	if err != nil {
		logrus.Errorf("svc.UserGet: set cache failed: %v", err)
		return nil, err
	}
	return userGetResp, nil
}

// setCacheTokenForUser 保存用户登录token与用户名的关系
func (svc UserService) setCacheTokenForUser(ctx context.Context, token, username string) {
	key := keyUserToken(token)
	svc.cache.Set(ctx, key, username, 30*time.Minute)
	logrus.Debugf("cache user login token:%s -> %s", username, key)
}

// queryCacheTokenForUser 保存用户登录token与用户名的关系
func (svc UserService) queryCacheTokenForUser(ctx context.Context, token string) (string, error) {
	cacheSessionKey := global.SessionIDCachePrefix + token
	username, err := svc.cache.Get(ctx, cacheSessionKey).Result()
	if err != nil {
		logrus.Errorf("queryCacheTokenForUser cache.Get error %v", err)
		return "", err
	}
	return username, err
}

func keyUserToken(token string) string {
	return global.SessionIDCachePrefix + token
}

//func generateToken(username string) string {
//	timestamp := time.Now().Format(time.RFC3339)
//	mux := timestamp + strconv.Itoa(rand.Int()) + username
//	h := md5.New()
//	h.Write([]byte(mux))
//	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
//}

package tcp_service

import (
	"context"
	"fmt"
	"simple_grpc/cmd/global"
	"simple_grpc/proto"
	"testing"
)

func TestUserService_Register(t *testing.T) {
	fmt.Println("Finished init logger...")
	global.DBEngine = InitDB()
	fmt.Println("Finished init mysql...")
	global.CacheClient = InitRedis()
	fmt.Println("Finished init redis...")
	svc := NewUserService(context.Background())

	// Mock stuffs
	username := "test_username2"
	nickname := "test_nickname2"
	password := "test_password2"

	// Input
	request := &proto.RegisterRequest{
		Username: username,
		Nickname: nickname,
		Password: password,
		//ProfilePic: "",
	}

	t.Run("normal register", func(t *testing.T) {
		// Test and compare with reflect.DeepEqual
		_, err := svc.Register(context.Background(), request)
		if err != nil {
			t.Errorf("TestUserService_Register got error %v", err)
		}
	})

	t.Run("invalid register", func(t *testing.T) {

		// should return an err
		_, err := svc.Register(context.Background(), request)
		if err == nil {
			t.Error("TestUserService_Register should return error but didn't")
		}
	})

}

func TestUserService_Login(t *testing.T) {
	fmt.Println("Finished init logger...")
	global.DBEngine = InitDB()
	fmt.Println("Finished init mysql...")
	global.CacheClient = InitRedis()
	fmt.Println("Finished init redis...")
	svc := NewUserService(context.Background())

	// Mock stuffs
	username := "test_username2"
	password := "test_password2"

	// Input
	request := &proto.LoginRequest{
		Username: username,
		Password: password,
	}

	t.Run("normal login", func(t *testing.T) {
		// Test and compare with reflect.DeepEqual
		_, err := svc.Login(context.Background(), request)
		if err != nil {
			t.Errorf("TestUserService_Login got error %v", err)
		}
	})

	t.Run("invalid login", func(t *testing.T) {
		// Test and compare with reflect.DeepEqual
		request.Password = "invalid"
		_, err := svc.Login(context.Background(), request)
		if err == nil {
			t.Errorf("TestUserService_Login got error %v", err)
		}
	})
}

func TestUserService_GetUser(t *testing.T) {
	fmt.Println("Finished init logger...")
	global.DBEngine = InitDB()
	fmt.Println("Finished init mysql...")
	global.CacheClient = InitRedis()
	fmt.Println("Finished init redis...")
	svc := NewUserService(context.Background())

	//// Mock stuffs
	//username := "test_username2"
	//password := "test_password2"

	// Input
	//request := &proto.LoginRequest{
	//	Username: username,
	//	Password: password,
	//}
	t.Run("normal login", func(t *testing.T) {
		// Test and compare with reflect.DeepEqual
		//SessionIdReply, err := svc.Login(context.Background(), request)
		//if err != nil {
		//	t.Errorf("TestUserService_Login got error %v", err)
		//}
		//fmt.Println("SessionIdReply:", SessionIdReply.SessionId)
		getUserrequest := &proto.GetUserRequest{
			//SessionId: SessionIdReply.SessionId,
			SessionId: "4e984a64-60f1-41c7-a442-ed12f05ed88b",
		}
		userInfo, err := svc.GetUser(context.Background(), getUserrequest)
		if err != nil {
			t.Errorf("TestUserService_GetUser got error %v", err)
		}
		fmt.Println("userInfo:", userInfo)
	})
}

func TestUserService_EditUser(t *testing.T) {
	fmt.Println("Finished init logger...")
	global.DBEngine = InitDB()
	fmt.Println("Finished init mysql...")
	global.CacheClient = InitRedis()
	fmt.Println("Finished init redis...")
	svc := NewUserService(context.Background())

	// Mock stuffs
	username := "test_username2"
	password := "test_password2"

	// Input
	request := &proto.LoginRequest{
		Username: username,
		Password: password,
	}
	t.Run("normal login", func(t *testing.T) {
		// Test and compare with reflect.DeepEqual
		SessionIdReply, err := svc.Login(context.Background(), request)
		if err != nil {
			t.Errorf("TestUserService_Login got error %v", err)
		}
		getUserrequest := &proto.GetUserRequest{
			SessionId: SessionIdReply.SessionId,
		}
		fmt.Println("SessionIdReply:", SessionIdReply.SessionId)
		userInfo, err := svc.GetUser(context.Background(), getUserrequest)
		if err != nil {
			t.Errorf("TestUserService_GetUser got error %v", err)
		}
		fmt.Println("userInfo:", userInfo)
		editUserrequest := &proto.EditUserRequest{
			SessionId:  SessionIdReply.SessionId,
			Nickname:   "update_Nickname",
			ProfilePic: "update_profile_pic",
		}
		editUserreply, err := svc.EditUser(context.Background(), editUserrequest)
		if err != nil {
			t.Errorf("TestUserService_GetUser got error %v", err)
		}
		fmt.Println("editUserreply:", editUserreply)
		updateUserInfo, err := svc.GetUser(context.Background(), getUserrequest)
		if err != nil {
			t.Errorf("TestUserService_GetUser got error %v", err)
		}
		fmt.Println("update userInfo:", updateUserInfo)
	})
}

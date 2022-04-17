package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"strings"
	"time"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"

	"shop_srvs/user_srv/global"
	"shop_srvs/user_srv/model"
	"shop_srvs/user_srv/proto"
)
/**
对密码在加以解释，md5+盐值是不可逆的(没有解密这一说)
但是加密和解密是有固定的规则的，通过一个options规则
 */
type UserServer struct {

}

// 获取用户列表
func (this *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	response := &proto.UserListResponse{}
	response.Total = int32(result.RowsAffected)

	// 取出数据并且进行了分页
	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)

	for _, user := range users {
		userInfoRsp := ModelToResponse(user)
		response.Data = append(response.Data, &userInfoRsp)
	}
	return response, nil
}

// 通过mobile查询user(手机号码查询用户)
func (this *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	userInfoResponse := ModelToResponse(user)
	return &userInfoResponse, nil
}

// 通过id查询user
func (this *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	userInfoResponse := ModelToResponse(user)
	return &userInfoResponse, nil
}

/*
新建用户
(1) check 用户是否存在（通过手机号check）
(2) 创建用户
 */
func (this *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已经存在")
	}

	user.Mobile = req.Mobile
	user.NickName = req.NickName
	user.Password = passwordEncryption(req.PassWord)

	result = global.DB.Create(&user)
	if result.Error != nil {  // Internal内部错误
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}

/*
个人用户修改信息
(1) 先查询用户是否存在
(2) 修改用户信息
 */
func (this *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {  // 用户不存在
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	// uint64 —> time
	birthDay := time.Unix(int64(req.BirthDay), 0)
	user.NickName = req.NickName
	user.Birthday = &birthDay
	user.Gender = req.Gender

	result = global.DB.Save(&user)  // 还可以用update
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &empty.Empty{}, nil
}

/*
校验密码（login业务需要使用）
 */
func (this *UserServer) CheckPassWord(ctx context.Context, req *proto.PassWordCheckInfo) (*proto.CheckResponse, error) {
	options := &password.Options{16, 100, 32, sha512.New}
	passwordInfo := strings.Split(req.EncryptedPassWord, "$")  // 加密密码
	check := password.Verify(req.PassWord, passwordInfo[2], passwordInfo[3], options)  // 真实密码
	return &proto.CheckResponse{Success: check}, nil
}

// 分页功能（参照：https://gorm.io/zh_CN/docs/scopes.html）
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// 将model 转化成 Response(userInfo)
func ModelToResponse(user model.User) proto.UserInfoResponse {
	// 在grpc的message中字段有默认值, 不能随便赋值nil进去, 在序列化会出错
	userInfoResponse := proto.UserInfoResponse{
		Id:       user.ID,
		PassWord: user.Password,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}

	// 不能将nil赋值进去
	if user.Birthday != nil {
		userInfoResponse.BirthDay = uint64(user.Birthday.Unix())
	}
	return userInfoResponse
}

// 密码加密（md5盐值加密）, 这种加密是不可逆的，只能通过加密密码和真实密码进行比较是否相同，check函数在上面
func passwordEncryption(pwd string) string {
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(pwd, options)
	// pbkdf2包名-sha512算法
	return fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
}



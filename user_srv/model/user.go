package model

import (
	"time"

	"gorm.io/gorm"
)

// 基础模型(每个变量名字都需要严格按照规范来)
type BaseModel struct {
	ID int32 `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"column:add_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`

	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

/*
用户表
密码: 1. 需要密文保存 2. 密文不可反解
     (1) 对称加密
	 (2) 非对称加密
	 (3) md5 信息摘要算法
	 密码如果不可反解, 用户找回密码。给用户一个链接, 进行验证个人信息修改密码
 */
type User struct {
	BaseModel
	Mobile string `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	Password string `gorm:"type:varchar(100);not null"`
	NickName string `gorm:"type:varchar(20)"`
	Birthday *time.Time `gorm:"type:datetime"`
	// 性别
	Gender string `gorm:"column:gender;default:male;type:varchar(6) comment 'female表示女, male表示男'"`
	// 区别是否是管理员
	Role int `gorm:"column:role;default:1;type:int comment '1表示普通用户, 2表示管理员'"`
}

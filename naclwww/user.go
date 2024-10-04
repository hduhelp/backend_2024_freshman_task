package main

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id       int
	Name     string
	Password string
	Cookie   string
}

func AddUser(u *User) error {
	//dsn := "hdubbs:yuhui123321@tcp(naclwww.xyz:3306)/HDUBBS?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//
	//if err != nil {
	//	return err
	//}
	//
	//_ = db.AutoMigrate(&User{})
	result := db.Where("name = ?", u.Name).First(u)
	if result.RowsAffected != 0 {
		return errors.New("user already exists")
	}
	result = db.Create(u)
	err := result.Error
	if err != nil {
		return errors.New("create user fail")
	} else {
		return nil
	}
}

func DeleteUser(u *User) error {
	if u.Id == 0 {
		return errors.New("no id")
	}
	if u.Name == "admin" {
		return errors.New("can not move admin")
	}
	result := db.Delete(u)
	err := result.Error
	if err != nil {
		panic("delete user fail,please check!")
	}
	return nil
}

func CheckLogin(u *User) (bool, error) {
	var user User
	result := db.Where("name = ? AND password= ?", u.Name, u.Password).First(&user)
	err := result.Error
	if err != nil || (result.RowsAffected) == 0 {
		return false, err
	} else {
		return true, nil
	}

	//if user.Password == u.Password {
	//	fmt.Println("pass")
	//	return true
	//} else {
	//	fmt.Println("fail")
	//	return false
	//}

}

func CookieInit(u *User) error { //ensure have user!!!
	currentTime := time.Now()
	u.Cookie = uuid.New().String() + currentTime.String()[:10]
	result := db.Model(&User{}).Where("name = ?", u.Name).Update("cookie", u.Cookie)

	err := result.Error
	if err != nil {
		db.Model(&User{}).Where("name = ?", u.Name).Update("cookie", "") //in case write cookie success
		return err
	} else {
		return nil
	}
}

func CookieCheck(u *User) (bool, error) { //check cookie exist and haven't out of data
	if u.Cookie == "" {
		return false, errors.New("cookie is empty,please login first")
	}

	var user User
	result := db.Where("cookie = ?", u.Cookie).First(&user)

	err := result.Error
	if err != nil || result.RowsAffected == 0 { // check exist
		return false, errors.New("illegal cookie")
	} else { //check time
		u.Name, u.Id = user.Name, user.Id
		t := u.Cookie[36:]
		Timetemplate := "2006-01-02"
		stamp, _ := time.ParseInLocation(Timetemplate, t, time.Local)
		stamp = stamp.Add(24 * 10 * time.Hour) // 10 days
		if stamp.Before(time.Now()) {
			return false, errors.New("cookie out of data")
		}
		return true, nil
	}
}

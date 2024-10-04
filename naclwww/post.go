package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Post struct {
	Id       int
	Title    string `default:""`
	Text     string
	Belongs  int
	Father   int
	Children string
	Time     time.Time
}

func AppendPost(post *Post) error {
	if post.Father != 0 {
		var fatherPost Post
		result := db.Where("id = ?", post.Father).First(&fatherPost)
		err := result.Error
		if err != nil {
			return err
		}

		result = db.Create(&post)
		err = result.Error
		if err != nil {
			return err
		}

		childrenList := fatherPost.Children + strconv.Itoa(post.Id) + " "
		fmt.Println(strconv.Itoa(post.Id))
		fmt.Println(childrenList)
		result = db.Model(&fatherPost).Update("children", childrenList)
		err = result.Error
		if err != nil {
			return err
		}
	} else {
		result := db.Create(&post)
		err := result.Error
		if err != nil {
			return err
		}
	}
	return nil
}

func ViewPost(id int) (*Post, error) {
	if id == 0 {
		return nil, errors.New("")
	}
	var post Post
	result := db.Where("id = ?", id).First(&post)
	err := result.Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func DeletePost(p *Post) error {
	if p.Id == 0 {
		return errors.New("no id")
	}
	result := db.Delete(p)
	err := result.Error
	if err != nil {
		panic("delete user fail,please check!")
	}
	return nil
}

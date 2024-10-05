package utils

import (
	"bytes"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
)

func UploadFile(fileBytes []byte, objectName string) string {
	client, err := oss.New("", "", "")
	if err != nil {
		log.Printf("创建 OSS 客户端失败: %v", err)
		return ""
	}
	bucket, err := client.Bucket("my-bilibili-project")
	if err != nil {
		log.Printf("获取存储空间失败: %v", err)
		return ""
	}
	err = bucket.PutObject(objectName, bytes.NewReader(fileBytes))
	if err != nil {
		log.Printf("上传文件失败: %v", err)
		return ""
	}
	return "https://my-bilibili-project.oss-cn-hangzhou.aliyuncs.com/" + objectName
}

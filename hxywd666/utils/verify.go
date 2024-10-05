package utils

import (
	"net/url"
	"strings"
)

func ValidatorURL(str string) bool {
	str = strings.TrimSpace(str)
	if str == "" {
		return false
	}
	u, err := url.Parse(str)
	if err != nil {
		return false
	}
	if u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}

func GetImageExtensionFromBase64(base64Str string) string {
	// 提取 Base64 编码中的数据部分
	dataPart := base64Str[strings.IndexByte(base64Str, ',')+1:]

	// 常见的图片 Base64 数据开头
	imagePrefixes := map[string]string{
		"data:image/jpeg;base64": "jpg",
		"data:image/png;base64":  "png",
		"data:image/gif;base64":  "gif",
	}

	for prefix, ext := range imagePrefixes {
		if strings.HasPrefix(dataPart, prefix) {
			return ext
		}
	}

	return ""
}

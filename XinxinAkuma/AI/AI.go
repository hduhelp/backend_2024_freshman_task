package AI

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	xunfeiAIAPIUrl = "wss://spark-api.xf-yun.com/v3.5/chat"
	apiSecret      = "OTM2NGMxOWJjY2FkOGYwZTEyOTVjZGY2"
	apiKey         = "ad54d6374685da80a5f420297ab6af00"
	appId          = "cee63188"
)

// GenerateSum 通过WebSocket与AI模型交互以生成答案
func GenerateSum(question string, answers []string) (string, error) {
	d := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	// 握手并建立websocket连接
	conn, resp, err := d.Dial(assembleAuthUrl1(xunfeiAIAPIUrl, apiKey, apiSecret), nil)
	if err != nil {
		return "", fmt.Errorf("连接失败: %s, %v", readResp(resp), err)
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn) // 确保在函数结束时关闭连接

	// 将所有的答案用 | 符号连接起来
	joinedAnswers := strings.Join(answers, "| ")

	// 构造最终的提示词
	prompt := fmt.Sprintf("我会给你一个问题和一组用 | 符号分隔的答案，帮我总结一个完整的回答，不要带有自己的评论和分析。 问题: %s\n答案: %s", question, joinedAnswers)
	data := genParams1(appId, prompt)

	// 发送数据
	if err := conn.WriteJSON(data); err != nil {
		return "", fmt.Errorf("发送数据失败: %v", err)
	}

	var answer string

	// 获取返回的数据
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return "", fmt.Errorf("读取消息失败: %v", err)
		}

		var data map[string]interface{}
		if err := json.Unmarshal(msg, &data); err != nil {
			return "", fmt.Errorf("解析JSON失败: %v", err)
		}
		// 解析数据
		payload, ok := data["payload"].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("无效的payload格式")
		}
		choices, ok := payload["choices"].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("无效的choices格式")
		}
		header, ok := data["header"].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("无效的header格式")
		}
		code, ok := header["code"].(float64)
		if !ok || code != 0 {
			return "", fmt.Errorf("错误的响应代码: %v", data["payload"])
		}

		status, ok := choices["status"].(float64)
		if !ok {
			return "", fmt.Errorf("无效的status格式")
		}
		text, ok := choices["text"].([]interface{})
		if !ok {
			return "", fmt.Errorf("无效的text格式")
		}
		content, ok := text[0].(map[string]interface{})["content"].(string)
		if !ok {
			return "", fmt.Errorf("无效的content格式")
		}

		if status != 2 {
			answer += content
		} else {
			answer += content
			usage, ok := payload["usage"].(map[string]interface{})
			if ok {
				temp, ok := usage["text"].(map[string]interface{})
				if ok {
					totalTokens, ok := temp["total_tokens"].(float64)
					if ok {
						fmt.Println("total_tokens:", totalTokens)
					}
				}
			}
			break
		}
	}

	// 输出返回结果
	return answer, nil
}

// 生成参数
func genParams1(appid, question string) map[string]interface{} { // 根据实际情况修改返回的数据结构和字段名

	messages := []Message{
		{Role: "user", Content: question},
	}

	data := map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
		"header": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"app_id": appid, // 根据实际情况修改返回的数据结构和字段名
		},
		"parameter": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"chat": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
				"domain":      "general",    // 根据实际情况修改返回的数据结构和字段名
				"temperature": float64(0.8), // 根据实际情况修改返回的数据结构和字段名
				"top_k":       int64(6),     // 根据实际情况修改返回的数据结构和字段名
				"max_tokens":  int64(2048),  // 根据实际情况修改返回的数据结构和字段名
				"auditing":    "default",    // 根据实际情况修改返回的数据结构和字段名
			},
		},
		"payload": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"message": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
				"text": messages, // 根据实际情况修改返回的数据结构和字段名
			},
		},
	}
	return data // 根据实际情况修改返回的数据结构和字段名
}

// 创建鉴权url  apikey 即 hmac username
func assembleAuthUrl1(hosturl string, apiKey, apiSecret string) string {
	ul, err := url.Parse(hosturl)
	if err != nil {
		fmt.Println(err)
	}
	//签名时间
	date := time.Now().UTC().Format(time.RFC1123)
	//date = "Tue, 28 May 2019 09:10:42 MST"
	//参与签名的字段 host ,date, request-line
	signString := []string{"host: " + ul.Host, "date: " + date, "GET " + ul.Path + " HTTP/1.1"}
	//拼接签名字符串
	sgin := strings.Join(signString, "\n")
	// fmt.Println(sgin)
	//签名结果
	sha := HmacWithShaTobase64("hmac-sha256", sgin, apiSecret)
	// fmt.Println(sha)
	//构建请求参数 此时不需要urlencoding
	authUrl := fmt.Sprintf("hmac username=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"", apiKey,
		"hmac-sha256", "host date request-line", sha)
	//将请求参数使用base64编码
	authorization := base64.StdEncoding.EncodeToString([]byte(authUrl))

	v := url.Values{}
	v.Add("host", ul.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)
	//将编码后的字符串url encode后添加到url后面
	callurl := hosturl + "?" + v.Encode()
	return callurl
}

func HmacWithShaTobase64(algorithm, data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}

func readResp(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("code=%d,body=%s", resp.StatusCode, string(b))
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

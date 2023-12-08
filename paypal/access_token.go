package paypal

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GenerateToken 返回一个access token, 用于商家创建订单
func GenerateToken() string {
	// 设置 PayPal 的客户端ID和客户端密钥
	clientID := "AeMmhtIC6Azh6dBFuLgqEdg3-RXdJ9QYWILHpvmWtUsF01EoT3gZnRl-rC8DtTEUoKdOiHtbh21VkDLz"
	clientSecret := "EOnSXtZAb4APnDWkUiZKUFqEqQMV6_VFTLvjPnP-hUBW1wcJoEzsHS04RzSeM4Qjx1jdD96rOlcB5iYC"

	// 编码 CLIENT_ID:CLIENT_SECRET 为 Base64
	credentials := clientID + ":" + clientSecret
	credentialsBase64 := base64.StdEncoding.EncodeToString([]byte(credentials))

	// 构建请求体
	data := "grant_type=client_credentials"
	requestBody := bytes.NewBuffer([]byte(data))

	// 创建 HTTP 客户端
	client := &http.Client{Timeout: 10 * time.Second}

	// 创建 HTTP POST 请求
	req, err := http.NewRequest("POST", "https://api-m.sandbox.paypal.com/v1/oauth2/token", requestBody)
	if err != nil {
		fmt.Println("创建HTTP POST 请求出错")
		return ""
	}

	// 设置请求头，包括授权头
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+credentialsBase64)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求时发生错误:", err)
		return ""
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应时发生错误:", err)
		return ""
	}

	// 创建一个结构体存储json 数据
	var responsemap map[string]interface{}
	if err := json.Unmarshal(body, &responsemap); err != nil {
		fmt.Println("解析出错")
		return ""
	}

	accessToken := responsemap["access_token"]
	return accessToken.(string)
}

package paypal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

// 创建一个结构来匹配 JSON 响应的格式
type ResponseData struct {
	ID            string `json:"id"`
	Status        string `json:"status"`
	PaymentSource struct {
		Paypal struct{} `json:"paypal"`
	} `json:"payment_source"`
	Links []struct {
		Href   string `json:"href"`
		Rel    string `json:"rel"`
		Method string `json:"method"`
	} `json:"links"`
}

// CreateOrder
//
// 需要商家生成一个uuid 作为paypalRequestID, 并且附带上订单的总金额以及access token;
//
// 返回一个map, 包含 订单ID, self, and payer-action 链接给用户，用户需要使用payer-action link 来执行支付
func CreateOrder(accessToken, paypalRequestID, amount string) map[string]string {
	// 生成 reference_id 和 PayPal-Request-Id
	referenceID := UUID()

	// 构建请求体
	requestBody := []byte(`{
		"intent": "CAPTURE",
		"purchase_units": [
			{
				"reference_id": "` + referenceID + `",
				"amount": {
					"currency_code": "USD",
					"value":"` + amount + `"
				},
				"shipping": {
					"address": {
						"address_line_1": "2211 N First Street",
						"address_line_2": "Building 17",
						"admin_area_2": "San Jose",
						"admin_area_1": "CA",
						"postal_code": "95131",
						"country_code": "US"
					}
				}
			}
		],
		"payment_source": {
			"paypal": {
				"experience_context": {
					"payment_method_preference": "IMMEDIATE_PAYMENT_REQUIRED",
					"brand_name": "EXAMPLE INC",
					"locale": "en-US",
					"landing_page": "LOGIN",
					"shipping_preference": "SET_PROVIDED_ADDRESS",
					"user_action": "PAY_NOW",
					"return_url": "https://example.com/returnUrl",
					"cancel_url": "https://example.com/cancelUrl"
				}
			}
		}
	}`)

	// 创建 HTTP 请求
	url := "https://api-m.sandbox.paypal.com/v2/checkout/orders"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("创建请求时发生错误:", err)
		return nil
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("PayPal-Request-Id", paypalRequestID)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求时发生错误:", err)
		return nil
	}
	defer resp.Body.Close()

	// 读取响应
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应时发生错误:", err)
		return nil
	}

	// 打印响应
	// fmt.Println("HTTP 响应码:", resp.Status)
	// fmt.Println("响应内容:", string(responseBody))

	responseData := new(ResponseData)
	resMap := make(map[string]string)
	if err := json.Unmarshal(responseBody, responseData); err != nil {
		fmt.Println("反序列化失败")
		return nil
	}

	// 存储订单ID
	resMap["id"] = responseData.ID

	// 遍历list
	for _, v := range responseData.Links {
		resMap[v.Rel] = v.Href
	}

	return resMap
}

// CapturePayment 发生在
func CapturePayment(orderId, paypalRequestId, accessToken string) (bool, error) {
	// PayPal API URL
	apiUrl := "https://api-m.sandbox.paypal.com/v2/checkout/orders/" + orderId + "/capture"

	// 构建请求体
	requestBody := []byte(`{}`) // 这里可以根据需要添加其他参数

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return false, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("PayPal-Request-Id", paypalRequestId)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode == http.StatusOK {
		// 订单支付捕获成功
		return true, nil
	} else {
		// 订单支付捕获失败
		return false, fmt.Errorf("Order capture failed with status code: %d", resp.StatusCode)
	}
}

func UUID() string {
	// 生成一个随机的UUID作为请求ID
	requestID := uuid.New().String()
	return requestID
}

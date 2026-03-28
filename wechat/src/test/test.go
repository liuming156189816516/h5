package main

import (
	"fmt"
	"net/url"
)

func main() {
	// 示例URL，包含查询参数
	rawURL := "https://b.xehh.cn/s12311/ua59h06dh1?3#/index"

	// 解析URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		panic(err) // 处理解析错误
	}

	// 获取查询参数
	queryParams := parsedURL.Query()

	// 提取特定参数，例如key2
	value := queryParams.Get("ttclid")     // 使用Get方法获取key2的值
	fmt.Println("Value of ttclid:", value) // 输出: Value of key2: value2
}

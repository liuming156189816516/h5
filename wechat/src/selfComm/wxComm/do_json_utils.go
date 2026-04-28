package wxComm

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 👉 返回结构
type AccountJson struct {
	Phone string      `json:"phone"`
	Token interface{} `json:"token"`
}

// 👉 递归处理 zip（支持嵌套 zip）
func processZip(r *zip.Reader) []AccountJson {

	var result []AccountJson

	for _, f := range r.File {

		if f.FileInfo().IsDir() {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			fmt.Println("打开失败:", f.Name)
			continue
		}

		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			fmt.Println("读取失败:", f.Name)
			continue
		}

		name := strings.ToLower(f.Name)

		// ✅ 处理 JSON 文件
		if strings.HasSuffix(name, ".json") {

			// 👉 解析 JSON（token）
			var obj interface{}
			if err := json.Unmarshal(data, &obj); err != nil {
				fmt.Println("非法JSON:", f.Name)
				continue
			}

			// 👉 文件名 → phone（只取 "_" 前面部分）
			baseName := strings.TrimSuffix(filepath.Base(f.Name), filepath.Ext(f.Name))
			parts := strings.Split(baseName, "_")

			phone := baseName
			if len(parts) > 0 {
				phone = parts[0]
			}

			result = append(result, AccountJson{
				Phone: phone,
				Token: obj,
			})

			continue
		}

		// ✅ 处理嵌套 zip
		if strings.HasSuffix(name, ".zip") {
			reader := bytes.NewReader(data)

			nestedZip, err := zip.NewReader(reader, int64(len(data)))
			if err != nil {
				fmt.Println("嵌套zip解析失败:", f.Name)
				continue
			}

			nested := processZip(nestedZip)
			result = append(result, nested...)
		}
	}

	return result
}

// 👉 主入口：传入目录，返回 Account 集合
func DoJsonUtils(dir string) ([]AccountJson, error) {

	var zipPaths []string

	// 👉 扫描目录
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".zip") {
			zipPaths = append(zipPaths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(zipPaths) == 0 {
		return nil, fmt.Errorf("目录下没有找到 zip 文件")
	}

	var allAccounts []AccountJson

	// 👉 处理所有 zip
	for _, path := range zipPaths {
		fmt.Println("处理:", path)

		r, err := zip.OpenReader(path)
		if err != nil {
			fmt.Println("打开失败:", path)
			continue
		}

		result := processZip(&r.Reader)
		allAccounts = append(allAccounts, result...)

		r.Close()
	}

	return allAccounts, nil
}

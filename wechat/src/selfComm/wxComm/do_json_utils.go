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

func processZip(r *zip.Reader, out *os.File) {
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

		if strings.HasSuffix(name, ".json") {
			var buf bytes.Buffer
			if err := json.Compact(&buf, data); err != nil {
				fmt.Println("非法 JSON:", f.Name)
				continue
			}

			out.Write(buf.Bytes())
			out.Write([]byte("\n"))
			continue
		}

		if strings.HasSuffix(name, ".zip") {
			reader := bytes.NewReader(data)
			nestedZip, err := zip.NewReader(reader, int64(len(data)))
			if err != nil {
				fmt.Println("嵌套zip解析失败:", f.Name)
				continue
			}
			processZip(nestedZip, out)
		}
	}
}

func DoJsonUtils(dir string) (string, error) {

	var zipNames []string
	var zipPaths []string

	// 👉 扫描指定目录
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".zip") {
			zipPaths = append(zipPaths, path)

			base := filepath.Base(path)
			name := strings.TrimSuffix(base, filepath.Ext(base))
			zipNames = append(zipNames, name)
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	if len(zipPaths) == 0 {
		return "", fmt.Errorf("目录下没有找到 zip 文件")
	}

	// 👉 输出文件（放在该目录下）
	outFileName := strings.Join(zipNames, "-") + ".txt"
	outPath := filepath.Join(dir, outFileName)

	out, err := os.Create(outPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	fmt.Println("输出文件:", outPath)

	// 👉 处理 zip
	for _, path := range zipPaths {
		fmt.Println("处理:", path)

		r, err := zip.OpenReader(path)
		if err != nil {
			fmt.Println("打开失败:", path)
			continue
		}

		processZip(&r.Reader, out)
		r.Close()
	}

	fmt.Println("完成 ✅")

	return outPath, nil
}

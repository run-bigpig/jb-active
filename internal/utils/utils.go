package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func SendRequest(url string, headers map[string]string, data interface{}) ([]byte, error) {
	var (
		req *http.Request
		err error
	)
	client := &http.Client{
		Timeout: time.Second * 50,
	}
	if data != nil {
		body, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
	}
	// 设置请求头（如果有）
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetCurrentPath 获取当前执行文件的路径
func GetCurrentPath() string {
	execPath, err := os.Executable()
	if err != nil {
		return ""
	}
	execPath = filepath.Dir(execPath)
	return execPath
}

// GetStaticPath 获取静态文件路径
func GetStaticPath() string {
	return filepath.Join(GetCurrentPath(), "static")
}

// IsExist 判断某个文件是否存在
func IsExist(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil || os.IsExist(err)
}

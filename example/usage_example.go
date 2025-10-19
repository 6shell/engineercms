package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// 使用示例：如何调用PDF扫描API

func main() {
	// 示例1：扫描目录
	scanDirectoryExample()

	// 示例2：查询单个文件信息
	getFileInfoExample()
}

// scanDirectoryExample 扫描目录示例
func scanDirectoryExample() {
	fmt.Println("=== PDF目录扫描示例 ===")

	// 构建请求
	reqBody := `{
		"directory": "E:/EngineerCMS6/attachment/standard",
		"max_workers": 5
	}`

	// 发送HTTP请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(
		"https://zsj.itdos.net/v1/wx/scan",
		"application/json",
		strings.NewReader(reqBody),
	)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 解析响应
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("响应解析失败: %v\n", err)
		return
	}

	// 格式化输出JSON结果
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("JSON格式化失败: %v\n", err)
		return
	}
	fmt.Printf("扫描结果:\n%s\n", string(jsonResult))
}

// getFileInfoExample 查询文件信息示例
func getFileInfoExample() {
	fmt.Println("\n=== 文件信息查询示例 ===")

	// 构建请求
	reqBody := `{
		"file_path": "J:/恢复的数据/D/temp2/CECS/CECS 193-2005城镇供水长距离输水管（渠）道工程技术规程.pdf"
	}`

	// 发送HTTP请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(
		"http://localhost:8080/api/v1/pdf/file-info",
		"application/json",
		strings.NewReader(reqBody),
	)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 解析响应
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("响应解析失败: %v\n", err)
		return
	}

	// 格式化输出JSON结果
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("JSON格式化失败: %v\n", err)
		return
	}
	fmt.Printf("文件信息:\n%s\n", string(jsonResult))
}

// 并发安全测试示例
func concurrentTestExample() {
	fmt.Println("\n=== 并发安全测试示例 ===")

	// 模拟多个并发请求
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			// 每个协程扫描不同的目录
			directory := fmt.Sprintf("/test/path/%d", index)
			reqBody := fmt.Sprintf(`{
				"directory": "%s",
				"max_workers": 2
			}`, directory)

			client := &http.Client{Timeout: 30 * time.Second}
			resp, err := client.Post(
				"http://localhost:8080/api/v1/pdf/scan",
				"application/json",
				strings.NewReader(reqBody),
			)
			if err != nil {
				fmt.Printf("协程 %d 请求失败: %v\n", index, err)
				return
			}
			defer resp.Body.Close()

			fmt.Printf("协程 %d 扫描完成\n", index)
		}(i)
	}

	wg.Wait()
	fmt.Println("所有并发扫描完成")
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Repository 结构体用于解析和编写 JSON 数据
type Repository struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	// 定义一个GET路由
	router.GET("/githubList", func(c *gin.Context) {
		// 打开文件
		file, err := os.Open("github_repos.json")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to open file: %s", err),
			})
			return
		}
		defer file.Close()

		// 直接将文件作为流发送
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Header().Set("Content-Type", "application/json")
		_, err = io.Copy(c.Writer, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to send file: %s", err),
			})
			return

		}
	})

	// 匹配/addRepo?owner=xxx&name=xxx
	router.GET("/addRepo", func(c *gin.Context) {
		owner := c.Query("owner")
		repoName := c.Query("name")
		// 检查请求参数是否完整
		if owner == "" || repoName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "owner and name parameters are required"})
			return
		}

		// 构造新的 repository 名称和 URL
		newRepo := Repository{
			Name: owner + "/" + repoName,
			URL:  "https://github.com/" + owner + "/" + repoName,
		}

		// 读取现有的 JSON 文件
		fileContents, err := ioutil.ReadFile("github_repos.json")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read json file"})
			return
		}

		// 解析 JSON 数据到数组
		var repos []Repository
		err = json.Unmarshal(fileContents, &repos)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse json file"})
			return
		}

		// 添加新的记录到数组
		repos = append(repos, newRepo)

		// 将更新后的数据编码回 JSON
		updatedJSON, err := json.Marshal(repos)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode json data"})
			return
		}

		// 将更新后的 JSON 写回文件
		err = ioutil.WriteFile("github_repos.json", updatedJSON, 0644)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write json file"})
			return
		}

		// 成功响应
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// 匹配/deleteRepo?url=xxx
	router.GET("/deleteRepo", func(c *gin.Context) {
		url := c.Query("url")
		// 检查请求参数是否完整
		// 读取现有 JSON 文件
		fileContents, err := ioutil.ReadFile("github_repos.json")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read json file"})
			return
		}

		// 解析 JSON 数据到数组
		var repos []Repository
		err = json.Unmarshal(fileContents, &repos)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse json file"})
			return
		}
		cnt := 0
		// 查找并删除指定的仓库
		newRepos := []Repository{}
		for _, repo := range repos {
			println("parse: ", url)
			println("repo: ", repo.URL)
			if url != repo.URL {
				newRepos = append(newRepos, repo)
				cnt++
			}
		}
		fmt.Println(cnt)
		// 将更新后的数据编码回 JSON
		updatedJSON, err := json.Marshal(newRepos)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode json data"})
			return
		}

		// 将更新后的 JSON 写回文件
		err = ioutil.WriteFile("github_repos.json", updatedJSON, 0644)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write json file"})
			return
		}

		// 成功响应
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	router.GET("/getIssues", func(c *gin.Context) {
		// 打开文件
		file, err := os.Open("github_issues" + time.Now().Format("20060102") + ".json")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to open file: %s", err),
			})
			return
		}
		defer file.Close()

		// 直接将文件作为流发送
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Header().Set("Content-Type", "application/json")
		_, err = io.Copy(c.Writer, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to send file: %s", err),
			})
			return

		}
	})
	// 在端口7853上启动服务器
	router.Run(":7853")
}

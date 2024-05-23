// 注意： 这个文件 是在生产环境启用的

//go:build deploy

package api

import "github.com/gin-gonic/gin"

func ginConfigure() {
	gin.SetMode(gin.ReleaseMode)
}

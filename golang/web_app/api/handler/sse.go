package handler

import (
	"fmt"
	"time"
	"github.com/gin-gonic/gin"
)

func SSE(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	fmt.Println(userId)
	// Content-Type を明示的に指定
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			// クライアントに送信
			fmt.Fprintf(c.Writer, "data: 現在時刻は %s\n\n", t.Format("15:04:05"))
			c.Writer.Flush() // バッファを明示的に送信
		case <-c.Request.Context().Done():
			// クライアント切断時
			fmt.Println("クライアント切断")
			return
		}
	}
}
package handler

import (
	// "fmt"
	"github.com/gin-gonic/gin"
	// "web_app/pubsub"
)

func SSE(c *gin.Context) {
	// // Content-Type を明示的に指定
	// c.Writer.Header().Set("Content-Type", "text/event-stream")
	// c.Writer.Header().Set("Cache-Control", "no-cache")
	// c.Writer.Header().Set("Connection", "keep-alive")

	// userId := c.MustGet("userId").(string)
	// fmt.Println(userId)

	// uuidSub, cancelSubscription := pubsub.SubscribeUuid(c.Request.Context(), userId)

	// for {			
	// 	select {
	// 	case msg := <-uuidSub:
	// 		// クライアントに送信
	// 		fmt.Fprintf(c.Writer, "data: UUID: %s\n\n", msg)
	// 		c.Writer.Flush() // バッファを明示的に送信
	// 	case <-c.Request.Context().Done():
	// 		// クライアント切断時
	// 		fmt.Println("クライアント切断")
	// 		cancelSubscription()
	// 		return
	// 	}
	// }
}
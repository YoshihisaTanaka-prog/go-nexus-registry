package pubsub

import (
	"context"
	"github.com/gin-gonic/gin"
	"sync"
)

type Hub struct {
	mu  sync.RWMutex
	bus map[string]*Bus // userId → Bus
}

var hub = &Hub{bus: make(map[string]*Bus)}

// GetBus は userId に対応する Bus を返す（なければ作成）
func getBus(userId string) *Bus {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	b, ok := hub.bus[userId]
	if ok {
		return b
	}
	b = &Bus{subs: make(map[string][]chan int)}
	hub.bus[userId] = b
	return b
}

// RemoveBus はログアウトなどで削除する際に使う
func RemoveBus(userId string) {
	hub.mu.Lock()
	delete(hub.bus, userId)
	hub.mu.Unlock()
}


type Bus struct {
	mu   sync.RWMutex
	subs map[string][]chan int // topic → subscribers
}

// subscribe はトピック購読チャネルを返す。ctx が終わると自動解除。
func (b *Bus) subscribe(ginContext context.Context, topic string) <-chan int {
	ctx, cancel := context.WithCancel(ginContext)

	ch := make(chan int, 32)

	b.mu.Lock()
	b.subs[topic] = append(b.subs[topic], ch)
	b.mu.Unlock()

	go func() {
		<-ctx.Done()
		b.mu.Lock()
		subs := b.subs[topic]
		for i, c := range subs {
			if c == ch {
				subs = append(subs[:i], subs[i+1:]...)
				break
			}
		}
		if len(subs) == 0 {
			delete(b.subs, topic)
		} else {
			b.subs[topic] = subs
		}
		close(ch)
		b.mu.Unlock()
	}()

	return ch
}
func Subscribe(ginContext *gin.Context, userId string, topic string) <-chan int {
	bus := getBus(userId)
	return bus.subscribe(ginContext.Request.Context(), topic)
}

// publish は非同期送信。
// 各購読チャネルのバッファが満杯ならその購読者には送らない（ドロップ）。
func (b *Bus) publish(topic string, msg int) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, ch := range b.subs[topic] {
		select {
		case ch <- msg:
		default: // drop
		}
	}
}

func Publish(userId string, topic string, msg int) {
	bus := getBus(userId)
	bus.publish(topic, msg)
}
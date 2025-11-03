package pubsub

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"sync"
)

type Hub struct {
	mu  sync.RWMutex
	bus map[string]*Bus[[2]int] // userId → Bus
}

var hub = &Hub{bus: make(map[string]*Bus[[2]int])}

type Bus[T any] struct {
	mu   sync.RWMutex
	subs map[string][]chan T // topic → subscribers
}

func getUuidString(uuId uuid.UUID) string {
	return fmt.Sprintf("%s", uuId)
}

// GetBus は userId に対応する Bus を返す（なければ作成）
func getBus(userId string) *Bus[[2]int] {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	b, ok := hub.bus[userId]
	if ok {
		return b
	}
	b = &Bus[[2]int]{subs: make(map[string][]chan [2]int)}
	hub.bus[userId] = b
	return b
}

// RemoveBus はログアウトなどで削除する際に使う
func RemoveBus(userId string) {
	hub.mu.Lock()
	delete(hub.bus, userId)
	hub.mu.Unlock()
}

var uuidBus = &Bus[uuid.UUID]{subs: make(map[string][]chan uuid.UUID)}

// subscribe はトピック購読チャネルを返す。ctx が終わると自動解除。
func (b *Bus[T]) subscribe(ginContext *context.Context, topic string) (<-chan T, context.CancelFunc) {
	ctx, cancel := context.WithCancel(*ginContext)

	ch := make(chan T, 32)

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

	return ch, cancel
}
func Subscribe(ginContext *context.Context, userId string, uuId uuid.UUID) (<-chan [2]int, context.CancelFunc) {
	bus := getBus(userId)
	return bus.subscribe(ginContext, getUuidString(uuId))
}

func SubscribeUuid(ginContext *context.Context, userId string) (<-chan uuid.UUID, context.CancelFunc) {
	return uuidBus.subscribe(ginContext, userId)
}

// publish は非同期送信。
// 各購読チャネルのバッファが満杯ならその購読者には送らない（ドロップ）。
func (b *Bus[T]) publish(topic string, msg T) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, ch := range b.subs[topic] {
		select {
		case ch <- msg:
		default: // drop
		}
	}
}

func Publish(userId string, uuId uuid.UUID, key int, status int) {
	bus := getBus(userId)
	bus.publish(getUuidString(uuId), [2]int{key, status})
}

func PublishUuid(userId string, msg uuid.UUID) {
	uuidBus.publish(userId, msg)
}
package ws

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"fortunefinger/backend/internal/pkg/errorx"

	"github.com/gorilla/websocket"
)

type PushMode uint8

const (
	PushAlways PushMode = iota + 1
	PushOnChange
)

const (
	defaultWriteTimeout = 10 * time.Second
	defaultReadTimeout  = 60 * time.Second
	defaultPingInterval = 25 * time.Second
)

type LoadFunc func(ctx context.Context) (any, error)
type FingerprintFunc func(data any) ([]byte, error)

type Topic struct {
	Key         string
	Interval    time.Duration
	Mode        PushMode
	Load        LoadFunc
	Fingerprint FingerprintFunc
}

type Message struct {
	Key  string `json:"key"`
	Data any    `json:"data"`
}

type Manager struct {
	mu       sync.RWMutex
	topics   map[string]Topic
	upgrader websocket.Upgrader
}

func NewManager() *Manager {
	return &Manager{
		topics: make(map[string]Topic),
		upgrader: websocket.Upgrader{
			HandshakeTimeout: 10 * time.Second,
			CheckOrigin:      checkOrigin,
		},
	}
}

func (m *Manager) MustRegister(topic Topic) {
	if topic.Key == "" || topic.Interval <= 0 || topic.Load == nil {
		panic("websocket topic configuration is invalid")
	}
	if topic.Mode != PushAlways && topic.Mode != PushOnChange {
		panic("websocket topic push mode is invalid")
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.topics[topic.Key]; exists {
		panic("websocket topic is already registered: " + topic.Key)
	}
	m.topics[topic.Key] = topic
}

func (m *Manager) Serve(ctx context.Context, w http.ResponseWriter, r *http.Request, key string) error {
	topic, exists := m.topic(key)
	if !exists {
		return errorx.New(errorx.NotFound, "WebSocket订阅主题不存在")
	}

	initial, err := topic.Load(ctx)
	if err != nil {
		return err
	}

	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil
	}
	defer conn.Close()

	conn.SetReadLimit(1024)
	_ = conn.SetReadDeadline(time.Now().Add(defaultReadTimeout))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(defaultReadTimeout))
	})

	done := make(chan struct{})
	go readLoop(conn, done)

	if err := writeMessage(conn, topic.Key, initial); err != nil {
		return nil
	}
	lastFingerprint, err := fingerprint(topic, initial)
	if err != nil {
		return nil
	}

	pushTicker := time.NewTicker(topic.Interval)
	pingTicker := time.NewTicker(defaultPingInterval)
	defer pushTicker.Stop()
	defer pingTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-done:
			return nil
		case <-pingTicker.C:
			if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(defaultWriteTimeout)); err != nil {
				return nil
			}
		case <-pushTicker.C:
			data, err := topic.Load(ctx)
			if err != nil {
				return nil
			}
			currentFingerprint, err := fingerprint(topic, data)
			if err != nil {
				return nil
			}
			if topic.Mode == PushOnChange && bytes.Equal(currentFingerprint, lastFingerprint) {
				continue
			}
			if err := writeMessage(conn, topic.Key, data); err != nil {
				return nil
			}
			lastFingerprint = currentFingerprint
		}
	}
}

func (m *Manager) topic(key string) (Topic, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	topic, exists := m.topics[key]
	return topic, exists
}

func readLoop(conn *websocket.Conn, done chan<- struct{}) {
	defer close(done)
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}

func writeMessage(conn *websocket.Conn, key string, data any) error {
	if err := conn.SetWriteDeadline(time.Now().Add(defaultWriteTimeout)); err != nil {
		return err
	}
	return conn.WriteJSON(Message{
		Key:  key,
		Data: data,
	})
}

func fingerprint(topic Topic, data any) ([]byte, error) {
	if topic.Fingerprint != nil {
		return topic.Fingerprint(data)
	}
	return json.Marshal(data)
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}
	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}
	originHost := parsed.Hostname()
	requestHost := r.Host
	if host, _, err := net.SplitHostPort(requestHost); err == nil {
		requestHost = host
	}
	return strings.EqualFold(originHost, requestHost)
}

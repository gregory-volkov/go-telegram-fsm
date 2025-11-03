package fsm

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StateID string

type Handler func(ctx *Context, update tgbotapi.Update) (StateID, error)

type Predicate func(update tgbotapi.Update) bool

type Transition struct {
	From StateID
	To   StateID
	Run  Handler
	If   Predicate
}

type SessionStore interface {
	Get(ctx context.Context, key string) (map[string]string, error)
	Set(ctx context.Context, key string, data map[string]string) error
}

type InMemoryStore struct {
	mu    sync.RWMutex
	store map[string]map[string]string
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{store: make(map[string]map[string]string)}
}

func (s *InMemoryStore) Get(ctx context.Context, key string) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if v, ok := s.store[key]; ok {
		copy := make(map[string]string, len(v))
		for k, vv := range v {
			copy[k] = vv
		}
		return copy, nil
	}
	return map[string]string{}, nil
}

func (s *InMemoryStore) Set(ctx context.Context, key string, data map[string]string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	copy := make(map[string]string, len(data))
	for k, v := range data {
		copy[k] = v
	}
	s.store[key] = copy
	return nil
}

type Context struct {
	Ctx    context.Context
	Bot    *tgbotapi.BotAPI
	Store  SessionStore
	Key    string
	Vars   map[string]string
	Logger func(format string, args ...interface{})
}

type FSM struct {
	bot         *tgbotapi.BotAPI
	store       SessionStore
	transitions map[StateID][]Transition
	initial     StateID
	mu          sync.RWMutex
	timeout     time.Duration
}

func NewFSM(bot *tgbotapi.BotAPI, store SessionStore, initial StateID, timeout time.Duration) *FSM {
	return &FSM{
		bot:         bot,
		store:       store,
		transitions: make(map[StateID][]Transition),
		initial:     initial,
		timeout:     timeout,
	}
}

func (f *FSM) RegisterTransition(t Transition) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.transitions[t.From] = append(f.transitions[t.From], t)
}

func sessionKeyFromUpdate(u tgbotapi.Update) string {
	if u.Message != nil {
		return "chat:" + strconv.FormatInt(u.Message.Chat.ID, 10)
	}
	if u.CallbackQuery != nil {
		return "chat:" + strconv.FormatInt(u.CallbackQuery.Message.Chat.ID, 10)
	}
	return "unknown"
}

func (f *FSM) ProcessUpdate(u tgbotapi.Update) error {
	key := sessionKeyFromUpdate(u)
	vars, err := f.store.Get(context.Background(), key)
	if err != nil {
		return err
	}
	current := StateID(vars["state"])
	if current == "" {
		current = f.initial
	}
	fmt.Println("StateID:", current)
	if u.CallbackQuery != nil {
		fmt.Println("Callback data:", u.CallbackQuery.Data)
	}

	f.mu.RLock()
	candidates := append([]Transition{}, f.transitions[current]...)
	f.mu.RUnlock()

	ctx := &Context{
		Ctx:    context.Background(),
		Bot:    f.bot,
		Store:  f.store,
		Key:    key,
		Vars:   vars,
		Logger: func(format string, args ...interface{}) {},
	}

	for _, t := range candidates {
		if t.If != nil && !t.If(u) {
			continue
		}
		next, err := t.Run(ctx, u)
		if err != nil {
			return err
		}
		if next != "" && string(next) != string(current) {
			ctx.Vars["state"] = string(next)
			if err := f.store.Set(context.Background(), key, ctx.Vars); err != nil {
				return err
			}
		}
		break
	}
	return nil
}

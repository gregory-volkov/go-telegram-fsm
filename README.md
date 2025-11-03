# Go Telegram FSM Framework

A lightweight **Go framework for building Telegram bots** using the **Finite State Machine (FSM)** pattern.

Instead of handling updates with endless `if/else` or `switch` blocks, you define **states**, **transitions**, and **handlers**.  
This lets you describe conversations as small, testable functions that handle user input and transition between states.

---

## 📂 Project structure

```

go-telegram-fsm/
├── fsm/
│   ├── fsm.go          # FSM engine and session handling
│   └── builder.go      # Fluent DSL for defining transitions
├── example/
│   └── main.go         # Demo bot
├── go.mod
├── go.sum
└── README.md

````

---

## 🚀 Getting Started

### 1. Clone & initialize

```bash
git clone https://github.com/gregory-volkov/go-telegram-fsm.git
cd go-telegram-fsm

go mod tidy
````

### 2. Set your bot token

Create a bot with [@BotFather](https://t.me/BotFather), then:

```bash
export TELEGRAM_TOKEN="123456789:ABCDEF..."
```

### 3. Run the example bot

```bash
go run example/main.go
```

### 4. Talk to your bot

* Send `/start`
* The bot asks: “What’s your name?”
* Reply with your name
* The bot greets you and returns to the start state

This demonstrates **stateful conversation flow** using an **in-memory session store**.

---

## 🛠️ Requirements

* Go 1.22 or later
* [Telegram Bot Token](https://t.me/BotFather)

---

## ⚖️ License

WTFPL -- do whatever you want with this.

---

## 💬 Author

Created by [Grigory Volkov](https://github.com/gregory-volkov) \
Contributions and issues welcome!

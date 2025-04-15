# Tasklight

**Tasklight** is a minimalist macOS app inspired by Spotlight, built to make task entry as fast and seamless as possible. Using a global hotkey, you can instantly open a lightweight input bar, type a natural language task like “Finish essay by Friday,” and have it automatically parsed with GPT and added to your Notion database.

---

## 🧠 Purpose

Tasklight was created to reduce the friction of capturing tasks. Rather than switching apps or losing focus, you can log tasks directly from anywhere on your system with just a keyboard shortcut. It’s perfect for fast-paced workflows and thought capture.

---

## ✨ Features

- Global hotkey to summon a Spotlight-style input window
- Transparent, always-on-top, distraction-free interface
- Natural language input processed through GPT (OpenAI)
- Automatically creates structured tasks in a Notion database
- Instant visibility toggle and keyboard-driven interaction

---

## ⚙️ Tech Stack

- **Wails v2** – Native macOS app framework (Go + Web)
- **Go** – Backend logic, Notion API, and hotkey registration
- **React** – Frontend UI
- **OpenAI GPT-4o** – Parses input into structured task data
- **Notion API** – Task storage and database integration
- **golang.design/x/hotkey** – Global hotkey support on macOS

---

## 💡 Inspiration

This project was inspired by [Coding With Lewis](https://youtu.be/lhjgj45x66Y?si=WroHyV6KREMvTNdW), who demonstrated a similar productivity concept. Tasklight builds on that foundation with added intelligence, Notion integration, and a refined user experience.

---

## 📦 Configuration

Tasklight uses an embedded `.env` file at build time for secrets and config values, including:

- `NOTION_DB_ID`
- `NOTION_SECRET`
- `OPENAI_API_KEY`

These are loaded and injected securely at runtime.

---

## 🚀 Usage

1. Press `Ctrl + Space` (or your configured shortcut) to launch the input window.
2. Type a task in natural language.
3. Hit Enter to send it to your Notion database.
4. Press `Escape` to hide the window.

---

## 🔭 Future Plans

Planned features and improvements include:

- ✅ Customizable hotkey settings
- ✅ Settings window to set secrets in UI
- ✅ Task history and recent entries
- ✅ Notion page selection or multiple database support
- ✅ Offline fallback with local queueing
- ✅ More advanced parsing (e.g. recurring tasks, tags)
- ✅ UI refinements and optional dark/light theme toggle

---

Made using Go, React, GPT, and a little obsession with clean interfaces.
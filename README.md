# LLM-Telegram-Chatbot

This Telegram Chatbot is a Golang-based bot that allows users to engage in conversations with a Language Model (LLM) using the telego library. This README provides an overview of the project and instructions on how to get started.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Usage](#usage)
- [Customizing the LLM](#customizing-the-llm)
- [License](#license)

## Prerequisites

Before using the Telegram Chatbot, ensure you have the following:

- Golang installed on your system.

## Getting Started

1. Clone the repository to your local machine:
```bash
git clone https://github.com/dxo1a/Telegram-LLM-Bot.git
```
2. Change directory to the project folder:
```bash
cd Telegram-LLM-Bot
```
4. Generate a Telegram Bot API key:
   - Visit the BotFather on Telegram to create your bot and obtain the API key.
   - Enter the key in [main.go](main.go) ``` botToken := "BOT_TOKEN_HERE" ```
  
## Usage

To start the bot, run the [main.go](main.go) file:
```bash
go run main.go
```
The bot doesn't have any commands, it's a basic sketch.

## Customizing the LLM

In the [main.go](main.go) file, you can customize the Language Model (LLM) used by the bot. You can also adjust LLM parameters and settings to suit your preferences.

## License

This project is licensed under the MIT License. Feel free to use and modify it as needed.

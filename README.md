# SSHMonitor

SSHMonitor is a lightweight program written in GoLang that monitors sshd login attempts and reports them to a Discord channel using a webhook. This tool is designed to enhance the security of your system by providing real-time notifications of SSH login attempts.

## Features

- Monitors sshd login attempts
- Sends notifications to a Discord channel via webhook

## Prerequisites

Before running SSHMonitor, make sure you have the following:

- GoLang installed on your system
- Discord account and a Discord webhook URL

## Setup

1. Clone the SSHMonitor repository to your local machine:

   ```bash
   git clone https://github.com/ExTBH/SSHMonitor.git
   ```

2. Navigate to the project directory:

   ```bash
   cd SSHMonitor
   ```

3. Create a `.env` file in the project root and specify your Discord webhook URL:

   ```plaintext
   WEBHOOK_URL=<your-discord-webhook-url>
   ```

## Usage

To run SSHMonitor, execute the following command:

```bash
go run SSHMonitor/SSHMonitor.go
```

SSHMonitor will start monitoring sshd login attempts and send notifications to the specified Discord channel when a login attempt is detected.
Logs will be saved to `./logfile.log` if no argument is passed, to specify a custom path 

```bash
go run SSHMonitor/SSHMonitor.go "/path/to/file.log"
```
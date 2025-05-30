# FFmate

FFmate is a modern and powerful automation layer built on top of FFmpeg - designed to make video and audio transcoding simpler, smarter, and easier to integrate.

## Overview

FFmate is an automation layer built on top of [FFmpeg](https://ffmpeg.org/), designed not only to simplify transcoding but also to serve as an extensible engine for custom media workflows. It provides developers with the tools to integrate FFmpeg's power into their applications and services through a comprehensive REST API, event-driven webhooks, and scriptable pre/post-processing hooks.

## Key Features

- **REST API** – Submit and manage FFmpeg tasks programmatically
- **Web UI** – Monitor and control jobs in real time, no terminal required
- **Watchfolders** – Automatically process files dropped into a directory
- **Presets** – Ready-made set of pre-configured transcoding presets for common use cases
- **Webhooks** – Get real-time notifications for task events
- **Dynamic Wildcards** – Automate file naming and folder structures
- **Pre/Post Processing** – Run custom scripts before or after each task to automate complex workflow steps
- **Built-in Queue** – Manage task execution with priority control and smart concurrency handling

## Documentation

Everything you need is available at [https://docs.ffmate.io](https://docs.ffmate.io)

## Prerequisites

- **FFmpeg**: Installed and in $PATH
- **Go (Golang)**: Version 1.24+ (Required for building from source or contributing to core)
- **Git**: For cloning the repository
- Familiarity with REST APIs and JSON
- Your preferred scripting language (Python, Bash, Node.js, etc.) for pre/post-processing scripts

## Installation

### Option 1: Use Pre-compiled Binaries

1. Download the latest binary from [GitHub Releases](https://github.com/welovemedia/ffmate/releases)
2. Make executable: `chmod +x ffmate`
3. (Optional) Move to a directory in your $PATH
4. Run FFmate Server: `./ffmate server`

### Option 2: Build from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/welovemedia/ffmate.git
   cd ffmate
   ```

2. Install Go dependencies:
   ```bash
   go mod tidy
   ```

3. Build the FFmate binary:
   ```bash
   go build -o ffmate main.go
   ```

4. Run FFmate Server:
   ```bash
   ./ffmate server
   ```

The API will be available at http://localhost:3000

## Usage Examples

- **Integrating with an existing application?**
  - Start with the `/api/v1/tasks` endpoint to submit jobs
  - Use `/api/v1/webhooks` to get status updates back to your application

- **Need custom logic around transcoding?**
  - Implement pre/post-processing scripts
  - Define `scriptPath` and `sidecarPath` in your task/preset API calls

- **Building a custom dashboard or UI?**
  - The entire Web UI is built on the public REST API and WebSockets
  - WebSockets (served from `/ws` on the same port) provide real-time updates for tasks and logs

## Contributing

We welcome contributions! See our contributing guidelines for more information.

## Community

- **Discord**: Join our [FFmate Community on Discord](https://discord.gg/NzfeHn37jT)
- **GitHub Issues**: Report bugs, discuss enhancements [GitHub Issues](https://github.com/welovemedia/ffmate/issues)

## License

This project is licensed under the AGPL-3.0. See the LICENSE file for full details.

# NetBird Slack Greeter

This Go-based application, residing in the `bot` folder, is designed to integrate with Slack for the NetBird community. It handles events from Slack's Socket Mode, offering automated responses to events of users joining the workspace or joining selected channels.

## Features

- **Automated Responses**: Automatically sends welcome messages to new users in specific channels or the entire workspace.
- **Template-based Messaging**: Employs predefined templates for different channels and events.
- **Event Handling**: Listens and responds to Slack events, like users joining a channel or the workspace.
- **Logging**: Uses `logrus` for structured logging, enhancing traceability and debugging.

## Requirements

- Go (tested with version 1.21.x or later, but probably work with older versions)
- Slack App Token (`SLACK_APP_TOKEN`) with the prefix `xapp-`
- Slack Bot Token (`SLACK_BOT_TOKEN`) with the prefix `xoxb-`

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/netbirdio/netbird-slack-greeter.git
   ```
2. Navigate to the project directory:
   ```bash
   cd netbird-slack-greeter
   ```
3. Build the project:
   ```bash
   cd bot
   go build -o netbird-slack-greeter
   ```

## Usage

1. Create a Slack app and install it to your workspace. Follow this [guide](https://api.slack.com/start/quickstart). You can use the manifest file from [manifest.yml](/manifest.yml).
2. Set the environment variables `SLACK_APP_TOKEN` and `SLACK_BOT_TOKEN`.
3. Update your constants and channel IDs in `bot/main.go`.
4. Build the project:
   ```bash
   cd bot
   go build -o netbird-slack-greeter
   ```
5. Execute the application:
   ```bash
   ./bot/netbird-slack-greeter
   ```

## Environment Variables

- `SLACK_APP_TOKEN`: Slack app-level token (must start with `xapp-`).
- `SLACK_BOT_TOKEN`: Slack bot token (must start with `xoxb-`).

## Logging

The project uses `logrus` for structured logging. It is configured to output logs in JSON format with timestamps and file locations for debugging.

## Message Templates

The application uses predefined templates for different types of messages:

- `bugsIssuesEtcTemplate`: For general queries and bug reporting.
- `selfHostedTemplate`: For self-hosted deployment discussions and support.
- `newUserTemplate`: For welcoming new users to the NetBird Slack community.

## Contributions

Contributions are encouraged! For feature requests or bug reports, please file an issue or submit a pull request on the GitHub repository.

## License

BSD 3-Clause License

Copyright (c) 2022 Wiretrustee UG (haftungsbeschränkt) & AUTHORS

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

## Support

For support, join our Slack community or file an issue on the GitHub repository.

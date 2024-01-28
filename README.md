# NetBird Slack Greeter

This Go-based application, residing in the `bot` folder, is designed to integrate with Slack for the NetBird community. It handles events from Slack's Socket Mode, offering automated responses and channel management.

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

1. Set the environment variables `SLACK_APP_TOKEN` and `SLACK_BOT_TOKEN`.
2. Update your constants and channel IDs in `bot/main.go`.
3. Build the project:
   ```bash
   cd bot
   go build -o netbird-slack-greeter
   ```
4. Execute the application:
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

[Specify License Here]

## Support

For support, join our Slack community or file an issue on the GitHub repository.

# GitLab MR Reviewer Bot

An AI-powered GitLab merge request reviewer bot that uses Google's Gemini LLM to provide intelligent code reviews. The bot automatically analyzes merge requests and posts constructive feedback as comments.

## Features

- 🤖 **AI-Powered Reviews**: Uses Google Gemini to analyze code changes and provide intelligent feedback
- 🔗 **GitLab Integration**: Seamless integration with GitLab webhooks
- 🚀 **Automatic Comments**: Posts review comments directly on merge requests
- 🔒 **Secure**: Supports webhook signature verification
- 📊 **Comprehensive Analysis**: Reviews code quality, security, performance, and best practices
- 🐳 **Containerized**: Ready-to-deploy Docker setup

## Prerequisites

- Go 1.21 or later
- GitLab access token with API permissions
- Google Gemini API key
- Docker (optional, for containerized deployment)

## Setup

### 1. Clone the Repository

```bash
git clone https://github.com/vinamra28/operator-reviewer.git
cd operator-reviewer
```

### 2. Configure Environment Variables

Copy the example environment file and fill in your credentials:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
GITLAB_TOKEN=your_gitlab_access_token_here
GITLAB_BASE_URL=https://gitlab.com
GEMINI_API_KEY=your_gemini_api_key_here
WEBHOOK_SECRET=your_webhook_secret_here
PORT=8080
```

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Run the Application

#### Option A: Direct Go Run
```bash
go run cmd/main.go
```

#### Option B: Docker Compose
```bash
docker-compose up --build
```

#### Option C: Docker Build
```bash
docker build -t gitlab-mr-reviewer .
docker run -p 8080:8080 --env-file .env gitlab-mr-reviewer
```

## GitLab Webhook Configuration

1. Go to your GitLab project/group settings
2. Navigate to **Webhooks**
3. Add a new webhook with:
   - **URL**: `http://your-server:8080/webhook`
   - **Secret Token**: Your `WEBHOOK_SECRET` value
   - **Trigger**: Select "Merge request events"
   - **SSL Verification**: Enable if using HTTPS

## How It Works

1. GitLab sends a webhook when a merge request is opened, reopened, or updated
2. The bot validates the webhook signature (if configured)
3. Fetches the merge request changes via GitLab API
4. Sends the code changes to Google Gemini for analysis
5. Posts AI-generated review comments back to the merge request

## API Endpoints

- `POST /webhook` - GitLab webhook endpoint
- `GET /health` - Health check endpoint

## Project Structure

```
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go          # Configuration management
│   ├── handlers/
│   │   └── webhook.go         # Webhook handlers
│   ├── models/
│   │   └── models.go          # Data structures
│   ├── server/
│   │   └── server.go          # HTTP server setup
│   └── services/
│       ├── gitlab.go          # GitLab API client
│       └── review.go          # Gemini AI integration
├── Dockerfile                 # Docker configuration
├── docker-compose.yml         # Docker Compose setup
├── go.mod                     # Go module definition
└── README.md                  # This file
```

## Required Tokens and Permissions

### GitLab Access Token
Create a GitLab access token with the following scopes:
- `api` - Full API access
- `read_api` - Read API access
- `read_repository` - Read repository access

### Gemini API Key
1. Go to [Google AI Studio](https://makersuite.google.com/app/apikey)
2. Create a new API key
3. Use this key as your `GEMINI_API_KEY`

## Security Considerations

- Always use HTTPS in production
- Set a strong `WEBHOOK_SECRET` for webhook verification
- Keep your GitLab token and Gemini API key secure
- Consider rate limiting for the webhook endpoint
- Run the application behind a reverse proxy (nginx, etc.)

## Troubleshooting

### Common Issues

1. **Webhook not triggering**: Check that the webhook URL is accessible and the secret matches
2. **API errors**: Verify your GitLab token has the required permissions
3. **Gemini errors**: Ensure your API key is valid and you have sufficient quota
4. **Connection issues**: Check network connectivity and firewall settings

### Logs

The application logs important events including:
- Webhook received events
- GitLab API calls
- Gemini API interactions
- Error conditions

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Support

For issues and questions:
1. Check the troubleshooting section
2. Open an issue on GitHub
3. Check GitLab and Gemini API documentation
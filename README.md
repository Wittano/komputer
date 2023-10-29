# komputer

Discord bot behave as like "komputer". One of character in Star Track parody series created by Dem3000

# Configuration

Preferred method to deploy for instance of bot is creating docker container using prepared Dockerfile image.
If you don't want to use docker, you can use `make install` to build and install program. It'll install
in `/opt/komputer` directory.

## Environment variable

Bot is configurable via a few environment variables:

- DISCORD_BOT_TOKEN= # Discord Bot Token (required)
- APPLICATION_ID= # Your bot application id (required)
- SERVER_GUID= # Your server ID, where bot will register commands (required)
- MONGODB_URI= # URI to connect MongoDB database (required)
- MONGODB_DB_NAME= # Name of main collection for komputer bot (required)
- RAPID_API_KEY= # API key for RapidAPIs. In project is use: HumorAP (optional)
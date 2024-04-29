# komputer

Discord bot behave as like "komputer". One of character in Star Track parody series created by Dem3000

# Installation
## Manual installation
Build project and installa with this simple command:
```sh
sudo make install
```
Bot installs in /opt/komputer directory
## For NixOS
You can use flake to install bot as NixOS module:

```nix
# Installa in flake
inputs = {
    komputer.url = "github:Wittano/komputer";
};

# Activate bot in NixOS
komputer = {
    enable = true;
    token = "<YOUR_DISCORD_TOKEN>";
    applicationID = "<YOUR_APPLICATION_ID>";
    guildID = "<YOUR_DISCORD_SERVER_ID>";
};
```

# Configuration

Preferred method to deploy for instance of bot is creating docker container using prepared Dockerfile image.
If you don't want to use docker, you can use `sudo make install` to build and install program. It'll install
in `/opt/komputer` directory.

## Environment variable

Bot is configurable via a few environment variables:

- DISCORD_BOT_TOKEN= Discord Bot Token (required)
- APPLICATION_ID= Your bot application id (required)
- SERVER_GUID= Your server ID, where bot will register commands (required)
- MONGODB_URI= URI to connect MongoDB database (required)
- RAPID_API_KEY= API key for RapidAPIs. In project is use: HumorAP (optional)


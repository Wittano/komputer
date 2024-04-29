# komputer

Discord bot behave as like "komputer". One of character in Star Track parody series created by Dem3000

# Commands

- /add-joke {category} {type} {answer} {question} - *komputer zapamiętaj dowcip* - Added new joke to database
    - category - REQUIRED - joke category
    - type - REQUIRED - joke type
    - answer - REQUIRED - joke content
    - question - OPTIONAL - additional question part of joke. Required only when you try add joke as **Two-Part** type
- /joke {id} {category} {type} - *komputer powiedz dowcip* - Find joke in JokeDev API, HumorAPI or your mongodb
  database. Each option in this command is optional:
    - id - jokeID, bot return ID after successful added joke into database
    - type - joke type
    - category - joke category
- /welcome - *komputer przywitaj się* - Bot welcome you
- spock {name} - *kurwa spock* - komputer says something
    - name - OPTIONAL - name of file, that bot should play. Name should be put without file extension e.g .mp3
- list {name} - *komputer co powiesz mi* - show list of available audio, that komputer play on voice chat
    - name - part of file name, which you want to see

## Joke categories

Bot recognizes a few categories of joke:

- Programowanie - programming joke
- Różne - any kind of joke, it's random
- YoMamma - special joke ;)

## Joke types

Bot recognized a two types of Joke:

- Single - Joke without question part. If you put joke in question field, this joke will be ignored
- Two-Parts - Joke with question part. You have to fill answer and question fields in /add-joke command

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
# Installa komputer via flake
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


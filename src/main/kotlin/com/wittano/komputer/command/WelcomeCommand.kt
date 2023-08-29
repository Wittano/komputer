package com.wittano.komputer.command

import discord4j.discordjson.json.ApplicationCommandRequest

class WelcomeCommand : Command {
    override fun execute() {
        TODO("Not yet implemented")
    }

    override fun createCommand(): ApplicationCommandRequest = ApplicationCommandRequest.builder()
        .name("welcome")
        .description("Welcome command to greetings to you")
        .build()
}
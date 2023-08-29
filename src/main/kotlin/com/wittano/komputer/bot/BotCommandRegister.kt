package com.wittano.komputer.bot

import com.wittano.komputer.command.AddJokeCommand
import com.wittano.komputer.command.JokeCommand
import com.wittano.komputer.command.WelcomeCommand
import com.wittano.komputer.config.ConfigLoader
import discord4j.rest.RestClient
import org.slf4j.LoggerFactory

class BotCommandRegister(private val client: RestClient) {

    private val logger = LoggerFactory.getLogger(BotCommandRegister::class.java)
    private val commands = arrayListOf(AddJokeCommand(), JokeCommand(), WelcomeCommand())

    fun singIn() {
        val config = ConfigLoader.load()

        commands.forEach {
            client.applicationService.createGuildApplicationCommand(
                config.applicationId,
                config.guildId,
                it.createCommand()
            ).doOnError { exception ->
                logger.error("Failed to register command. Cause: ${exception.message}")
            }.subscribe()
        }
    }

}
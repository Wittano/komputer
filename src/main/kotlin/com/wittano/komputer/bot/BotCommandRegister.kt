package com.wittano.komputer.bot

import com.google.inject.Inject
import com.wittano.komputer.command.SlashCommand
import com.wittano.komputer.config.ConfigLoader
import discord4j.rest.RestClient
import org.slf4j.LoggerFactory

class BotCommandRegister @Inject constructor(
    private val commands: Set<SlashCommand>,
) {
    private val logger = LoggerFactory.getLogger(BotCommandRegister::class.java)

    // TODO Replace my command register via discord4j register method: https://docs.discord4j.com/interactions/application-commands#simplifying-the-lifecycle
    fun singIn(client: RestClient) {
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
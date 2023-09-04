package com.wittano.komputer.core.bot

import com.wittano.komputer.core.config.ConfigLoader
import discord4j.core.DiscordClientBuilder
import discord4j.core.GatewayDiscordClient
import org.slf4j.LoggerFactory

val discordClient: GatewayDiscordClient by lazy {
    val log = LoggerFactory.getLogger(GatewayDiscordClient::class.qualifiedName)
    val config = ConfigLoader.load()

    DiscordClientBuilder.create(config.token)
        .build()
        .login()
        .doOnSuccess { log.info("Bot is ready!") }
        .block() ?: throw IllegalStateException("Failed to start up discord bot")
}
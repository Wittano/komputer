package com.wittano.komputer.core.bot

import com.wittano.komputer.core.config.config
import discord4j.core.DiscordClientBuilder
import discord4j.core.GatewayDiscordClient
import org.slf4j.LoggerFactory

val discordClient: GatewayDiscordClient by lazy {
    val log = LoggerFactory.getLogger(GatewayDiscordClient::class.qualifiedName)

    DiscordClientBuilder.create(config.token)
        .build()
        .login()
        .doOnSuccess { log.info("Bot is ready!") }
        .block() ?: throw IllegalStateException("Failed to start up discord bot")
}
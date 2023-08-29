package com.wittano.komputer.bot

import com.wittano.komputer.config.ConfigLoader
import discord4j.core.DiscordClientBuilder
import org.slf4j.LoggerFactory
import picocli.CommandLine.Command

@Command(
    name = "komputer",
    description = ["Discord bot behave as like \"komputer\". One of character in Star Track parody series created by Dem3000"]
)
class KomputerBot : Runnable {

    private val log = LoggerFactory.getLogger(this::class.qualifiedName)

    override fun run() {
        val config = ConfigLoader.load()
        val client = DiscordClientBuilder.create(config.token)
            .build()
            .login()
            .doOnSuccess { log.info("Bot is ready!") }
            .block()

        client?.onDisconnect()?.block()
    }

}
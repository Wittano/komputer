package com.wittano.komputer.bot

import com.google.inject.ConfigurationException
import com.google.inject.Injector
import com.google.inject.Key
import com.google.inject.name.Names
import com.wittano.komputer.command.SlashCommand
import com.wittano.komputer.config.ConfigLoader
import discord4j.core.DiscordClientBuilder
import discord4j.core.event.domain.interaction.ApplicationCommandInteractionEvent
import org.slf4j.LoggerFactory
import picocli.CommandLine.Command
import reactor.core.publisher.Mono

@Command(
    name = "komputer",
    description = ["Discord bot behave as like \"komputer\". One of character in Star Track parody series created by Dem3000"]
)
class KomputerBot(private val injector: Injector) : Runnable {

    private val log = LoggerFactory.getLogger(this::class.qualifiedName)

    override fun run() {
        val config = ConfigLoader.load()
        val client = DiscordClientBuilder.create(config.token)
            .build()
            .login()
            .doOnSuccess { log.info("Bot is ready!") }
            .block() ?: throw IllegalStateException("Failed to start up discord bot")

        BotCommandRegister(client.restClient).singIn()

        client.on(ApplicationCommandInteractionEvent::class.java) { event ->
            try {
                val slashCommand: SlashCommand =
                    injector.getInstance(
                        Key.get(
                            SlashCommand::class.java,
                            Names.named(event.commandName.replace("-", ""))
                        )
                    )

                return@on Mono.from<Nothing> {
                    slashCommand.execute(event)
                }
            } catch (ex: ConfigurationException) {
                return@on Mono.error(ex)
            }
        }.subscribe()

        client.onDisconnect().block()
    }

}
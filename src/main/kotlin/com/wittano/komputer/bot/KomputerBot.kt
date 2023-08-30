package com.wittano.komputer.bot

import com.google.inject.ConfigurationException
import com.google.inject.Injector
import com.google.inject.Key
import com.google.inject.name.Names
import com.wittano.komputer.command.SlashCommand
import com.wittano.komputer.config.ConfigLoader
import com.wittano.komputer.message.interaction.ButtonReaction
import discord4j.core.DiscordClientBuilder
import discord4j.core.event.domain.interaction.ButtonInteractionEvent
import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
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

        val commandRegister = injector.getInstance(BotCommandRegister::class.java)
        commandRegister.singIn(client.restClient)

        client.on(ChatInputInteractionEvent::class.java) {
            try {
                val slashCommand: SlashCommand = injector.getInstance(
                    Key.get(
                        SlashCommand::class.java,
                        Names.named(it.commandName.replace("-", ""))
                    )
                )

                return@on slashCommand.execute(it)
            } catch (ex: ConfigurationException) {
                return@on Mono.error(ex)
            }
        }.subscribe()

        client.on(ButtonInteractionEvent::class.java) {
            try {
                val buttonReaction = injector.getInstance(
                    Key.get(
                        ButtonReaction::class.java,
                        Names.named(it.customId.replace("-", ""))
                    )
                )

                return@on buttonReaction.execute(it)
            } catch (ex: ConfigurationException) {
                return@on Mono.error(ex)
            }
        }.subscribe()

        client.onDisconnect().block()
    }

}
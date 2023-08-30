package com.wittano.komputer.bot

import com.google.inject.ConfigurationException
import com.google.inject.Injector
import com.google.inject.Key
import com.google.inject.name.Names
import com.wittano.komputer.command.SlashCommand
import com.wittano.komputer.command.registred.RegisteredCommandsUtils
import com.wittano.komputer.config.ConfigLoader
import com.wittano.komputer.message.interaction.ButtonReaction
import discord4j.core.DiscordClientBuilder
import discord4j.core.GatewayDiscordClient
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

        val commands = RegisteredCommandsUtils.getCommandsFromJsonFiles()
        val registeredCommands =
            client.restClient.applicationService.getGuildApplicationCommands(config.applicationId, config.guildId)

        BotCommandCleaner.deleteUnusedGuildCommands(client.restClient, commands, registeredCommands)
            .thenMany(BotCommandRegister.registerCommands(client.restClient, commands, registeredCommands))
            .subscribe()

        handleChatInputEvents(client)
        handleButtonInteractionEvents(client)

        client.onDisconnect().block()
    }

    private fun handleButtonInteractionEvents(client: GatewayDiscordClient) {
        client.on(ButtonInteractionEvent::class.java) {
            try {
                val buttonReaction = injector.getInstance(
                    Key.get(
                        ButtonReaction::class.java,
                        Names.named(it.customId.replace("-", ""))
                    )
                )

                buttonReaction.execute(it)
            } catch (ex: ConfigurationException) {
                Mono.error(ex)
            }
        }.subscribe()
    }

    private fun handleChatInputEvents(client: GatewayDiscordClient) {
        client.on(ChatInputInteractionEvent::class.java) {
            try {
                val slashCommand: SlashCommand = injector.getInstance(
                    Key.get(
                        SlashCommand::class.java,
                        Names.named(it.commandName.replace("-", ""))
                    )
                )

                slashCommand.execute(it)
            } catch (ex: ConfigurationException) {
                Mono.error(ex)
            }
        }.subscribe()
    }

}
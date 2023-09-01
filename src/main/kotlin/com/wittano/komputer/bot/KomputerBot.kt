package com.wittano.komputer.bot

import com.wittano.komputer.command.registred.RegisteredCommandsUtils
import com.wittano.komputer.config.ConfigLoader
import com.wittano.komputer.config.dagger.DaggerKomputerComponent
import com.wittano.komputer.message.createErrorMessage
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
class KomputerBot : Runnable {

    private val log = LoggerFactory.getLogger(this::class.qualifiedName)
    private val komputerComponent = DaggerKomputerComponent.create()

    override fun run() {
        val config = ConfigLoader.load()
        val client = DiscordClientBuilder.create(config.token)
            .build()
            .login()
            .doOnSuccess { log.info("Bot is ready!") }
            .block() ?: throw IllegalStateException("Failed to start up discord bot")

        val commands = RegisteredCommandsUtils.getCommandsFromJsonFiles()
        val registeredCommands =
            client.restClient.applicationService
                .getGuildApplicationCommands(config.applicationId, config.guildId)
                .collectList()
                .filter {
                    it.isNotEmpty()
                }

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
                val customId = it.customId.replace("-", "")
                val buttonReaction = DaggerKomputerComponent.create().getButtonReaction()[customId]

                buttonReaction?.execute(it)
                    ?: Mono.error(NoSuchElementException("Button with id $customId wasn't found"))
            } catch (ex: Exception) {
                log.error("Unexpected error during handling '${it.customId}' button interaction", ex)
                it.reply(createErrorMessage())
            }
        }.subscribe()
    }

    private fun handleChatInputEvents(client: GatewayDiscordClient) {
        client.on(ChatInputInteractionEvent::class.java) {
            try {
                val commandName = it.commandName.replace("-", "")
                val slashCommand = komputerComponent.getSlashCommand()[commandName]

                slashCommand?.execute(it)
                    ?: Mono.error(NoSuchElementException("Slash command '$commandName' wasn't found"))
            } catch (ex: Exception) {
                log.error("Unexpected error during handling '${it.commandName}' chat interaction", ex)
                it.reply(createErrorMessage())
            }
        }.subscribe()
    }

}
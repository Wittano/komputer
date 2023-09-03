package com.wittano.komputer.core.bot

import com.wittano.komputer.config.dagger.DaggerKomputerComponent
import com.wittano.komputer.core.command.exception.CommandException
import com.wittano.komputer.core.command.registred.RegisteredCommandsUtils
import com.wittano.komputer.core.config.ConfigLoader
import com.wittano.komputer.core.joke.JokeException
import com.wittano.komputer.core.message.createErrorMessage
import com.wittano.komputer.message.resource.MessageResource
import discord4j.core.DiscordClientBuilder
import discord4j.core.GatewayDiscordClient
import discord4j.core.event.domain.interaction.ButtonInteractionEvent
import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import discord4j.core.event.domain.interaction.DeferrableInteractionEvent
import discord4j.core.spec.InteractionApplicationCommandCallbackSpec
import org.slf4j.LoggerFactory
import picocli.CommandLine.Command
import reactor.core.publisher.Mono
import java.util.*

// TODO Export cli option into new submodule e.g. cli
@Command(
    name = "komputer",
    description = ["Discord bot behave as like \"komputer\". One of character in Star Track parody series created by Dem3000"]
)
class KomputerBot : Runnable {

    private val log = LoggerFactory.getLogger(this::class.qualifiedName)
    private val komputerComponents = DaggerKomputerComponent.create()

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

        // TODO Export command registration management to separate subcommand
        BotCommandCleaner.deleteUnusedGuildCommands(client.restClient, commands, registeredCommands)
            .thenMany(BotCommandRegister.registerCommands(client.restClient, commands, registeredCommands))
            .subscribe()

        handleChatInputEvents(client)
        handleButtonInteractionEvents(client)

        client.onDisconnect().block()
    }

    private fun handleButtonInteractionEvents(client: GatewayDiscordClient) {
        client.on(ButtonInteractionEvent::class.java) { event ->
            val customId = event.customId.replace("-", "")
            val buttonReaction = komputerComponents.getButtonReaction()[customId]

            val errorResponse = Mono.error<Void>(CommandException("Button with id $customId wasn't found", customId))
                .doOnError { exception ->
                    val buttonIdError = exception.takeIf { it is CommandException }
                        ?.let { it as CommandException }
                        ?.let { "'${it.commandId}'" }
                        .orEmpty()

                    log.error("Unexpected error during handling $buttonIdError button interaction", exception)
                }.transform { event.reply(createErrorMessage()) }

            buttonReaction?.execute(event)
                ?.onErrorResume { exception -> sendErrorMessage(event, exception) }
                ?: errorResponse
        }.subscribe()
    }

    private fun handleChatInputEvents(client: GatewayDiscordClient) {
        client.on(ChatInputInteractionEvent::class.java) { event ->
            val commandName = event.commandName.replace("-", "")
            val slashCommand = komputerComponents.getSlashCommand()[commandName]

            val errorResponse =
                Mono.error<Void>(CommandException("Slash command '$commandName' wasn't found", commandName))
                    .doOnError { exception ->
                        val commandIdError = exception.takeIf { it is CommandException }
                            ?.let { it as CommandException }
                            ?.let { "'${it.commandId}'" }
                            .orEmpty()

                        log.error("Unexpected error during handling $commandIdError chat interaction", exception)
                    }.transform { event.reply(createErrorMessage()) }

            slashCommand?.execute(event)
                ?.onErrorResume { exception -> sendErrorMessage(event, exception) }
                ?: errorResponse
        }.subscribe()
    }

    private fun sendErrorMessage(
        event: DeferrableInteractionEvent,
        exception: Throwable,
        isUserOnlyVisible: Boolean = false
    ): Mono<Void> {
        val errorMessage = exception.takeIf { it is JokeException }
            ?.let { it as JokeException }
            ?.let {
                val locale = event.interaction.userLocale.split("-")
                    .takeIf { isUserOnlyVisible }
                    ?.let { (language, country) ->
                        Locale(language, country)
                    } ?: Locale("pl")
                val msg = MessageResource.get(it.code, locale)

                InteractionApplicationCommandCallbackSpec.builder()
                    .content(msg)
                    .build()
                    .withEphemeral(isUserOnlyVisible)
            }

        log.error("During execute command, was thrown unexpected error", exception)

        return event.reply(errorMessage ?: createErrorMessage())
    }

}
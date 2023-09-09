package com.wittano.komputer.bot.bot

import com.wittano.komputer.bot.command.exception.CommandException
import com.wittano.komputer.bot.dagger.DaggerKomputerComponent
import com.wittano.komputer.bot.joke.JokeException
import com.wittano.komputer.bot.message.createErrorMessage
import discord4j.core.event.domain.interaction.ButtonInteractionEvent
import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import discord4j.core.event.domain.interaction.DeferrableInteractionEvent
import org.slf4j.LoggerFactory
import reactor.core.publisher.Mono
import java.util.*

class KomputerBot {

    private val log = LoggerFactory.getLogger(this::class.qualifiedName)
    private val komputerComponents = DaggerKomputerComponent.create()

    fun start() {
        handleChatInputEvents()
        handleButtonInteractionEvents()

        discordClient.onDisconnect().block()
    }

    private fun handleButtonInteractionEvents() {
        discordClient.on(ButtonInteractionEvent::class.java) { event ->
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

    private fun handleChatInputEvents() {
        discordClient.on(ChatInputInteractionEvent::class.java) { event ->
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
    ): Mono<Void> {
        val errorMessage = exception.takeIf { it is JokeException }
            ?.let { it as JokeException }
            ?.let {
                // TODO Add global language in configuration
                val locale = event.interaction.userLocale.split("-")
                    .let { (language, country) ->
                        Locale(language, country)
                    }

                createErrorMessage(it.code, locale)
            }

        log.error("During execute command, was thrown unexpected error", exception)

        return event.reply(errorMessage ?: createErrorMessage())
    }

}
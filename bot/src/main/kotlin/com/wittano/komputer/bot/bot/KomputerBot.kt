package com.wittano.komputer.bot.bot

import com.wittano.komputer.bot.command.exception.CommandException
import com.wittano.komputer.bot.dagger.DaggerKomputerComponent
import com.wittano.komputer.bot.joke.JokeException
import com.wittano.komputer.bot.message.createErrorMessage
import com.wittano.komputer.commons.transtation.getErrorMessage
import discord4j.core.GatewayDiscordClient
import discord4j.core.event.domain.interaction.ButtonInteractionEvent
import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import discord4j.core.event.domain.interaction.DeferrableInteractionEvent
import discord4j.core.spec.InteractionApplicationCommandCallbackSpec
import org.slf4j.LoggerFactory
import reactor.core.publisher.Mono
import java.util.*

class KomputerBot {

    private val log = LoggerFactory.getLogger(this::class.qualifiedName)
    private val komputerComponents = DaggerKomputerComponent.create()

    fun start() {
        handleChatInputEvents(discordClient)
        handleButtonInteractionEvents(discordClient)

        discordClient.onDisconnect().block()
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
                val msg = getErrorMessage(it.code, locale)

                InteractionApplicationCommandCallbackSpec.builder()
                    .content(msg)
                    .build()
                    .withEphemeral(isUserOnlyVisible)
            }

        log.error("During execute command, was thrown unexpected error", exception)

        return event.reply(errorMessage ?: createErrorMessage())
    }

}
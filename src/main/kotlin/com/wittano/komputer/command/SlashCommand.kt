package com.wittano.komputer.command

import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import discord4j.discordjson.json.ApplicationCommandRequest
import reactor.core.publisher.Mono

interface SlashCommand {

    fun execute(event: ChatInputInteractionEvent): Mono<Void>

    fun createCommand(): ApplicationCommandRequest

}
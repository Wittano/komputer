package com.wittano.komputer.command

import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import reactor.core.publisher.Mono

fun interface SlashCommand {

    fun execute(event: ChatInputInteractionEvent): Mono<Void>

}
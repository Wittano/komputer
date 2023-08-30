package com.wittano.komputer.message.interaction

import discord4j.core.event.domain.interaction.ButtonInteractionEvent
import reactor.core.publisher.Mono

fun interface ButtonReaction {
    fun execute(event: ButtonInteractionEvent): Mono<Void>
}
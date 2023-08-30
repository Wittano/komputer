package com.wittano.komputer.command

import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import reactor.core.publisher.Mono

class AddJokeCommand : SlashCommand {
    override fun execute(event: ChatInputInteractionEvent): Mono<Void> {
        TODO("Not yet implemented")
    }
}
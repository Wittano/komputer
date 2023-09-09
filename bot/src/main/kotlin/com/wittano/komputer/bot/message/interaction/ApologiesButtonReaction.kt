package com.wittano.komputer.bot.message.interaction

import discord4j.core.event.domain.interaction.ButtonInteractionEvent
import reactor.core.publisher.Mono

class ApologiesButtonReaction : ButtonReaction {
    override fun execute(event: ButtonInteractionEvent): Mono<Void> {
        return event.reply("Przepraszam")
    }
}
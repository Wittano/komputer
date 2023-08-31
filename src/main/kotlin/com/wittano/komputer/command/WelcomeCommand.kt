package com.wittano.komputer.command

import com.wittano.komputer.utils.toNullable
import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import reactor.core.publisher.Mono
import javax.inject.Named

@Named("welcomeCommand")
class WelcomeCommand : SlashCommand {
    override fun execute(event: ChatInputInteractionEvent): Mono<Void> {
        return event.reply("Witaj kapitanie ${event.interaction.member.toNullable()?.displayName}")
    }
}
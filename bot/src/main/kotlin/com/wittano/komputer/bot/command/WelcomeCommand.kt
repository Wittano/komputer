package com.wittano.komputer.bot.command

import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import reactor.core.publisher.Mono
import javax.inject.Named
import kotlin.jvm.optionals.getOrNull

@Named("welcomeCommand")
class WelcomeCommand : SlashCommand {
    override fun execute(event: ChatInputInteractionEvent): Mono<Void> {
        return event.reply("Witaj kapitanie ${event.interaction.member.getOrNull()?.displayName}")
    }
}
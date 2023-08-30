package com.wittano.komputer.command

import com.wittano.komputer.utils.toNullable
import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import discord4j.discordjson.json.ApplicationCommandRequest
import jakarta.inject.Named
import reactor.core.publisher.Mono

@Named("welcomeCommand")
class WelcomeCommand : SlashCommand {
    override fun execute(event: ChatInputInteractionEvent): Mono<Void> {
        return event.reply("Witaj kapitanie ${event.interaction.member.toNullable()?.displayName}")
    }

    override fun createCommand(): ApplicationCommandRequest = ApplicationCommandRequest.builder()
        .name("welcome")
        .description("Welcome command to greetings to you")
        .build()
}
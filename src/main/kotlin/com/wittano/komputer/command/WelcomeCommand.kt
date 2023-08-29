package com.wittano.komputer.command

import com.wittano.komputer.utils.toNullable
import discord4j.core.event.domain.interaction.ApplicationCommandInteractionEvent
import discord4j.core.spec.InteractionApplicationCommandCallbackReplyMono
import discord4j.discordjson.json.ApplicationCommandRequest
import jakarta.inject.Named

@Named("welcomeCommand")
class WelcomeCommand : SlashCommand {
    override fun execute(event: ApplicationCommandInteractionEvent): InteractionApplicationCommandCallbackReplyMono {
        return event.reply("Witaj kapitanie ${event.interaction.member.toNullable()?.displayName}")
    }

    override fun createCommand(): ApplicationCommandRequest = ApplicationCommandRequest.builder()
        .name("welcome")
        .description("Welcome command to greetings to you")
        .build()
}
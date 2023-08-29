package com.wittano.komputer.command

import discord4j.core.event.domain.interaction.ApplicationCommandInteractionEvent
import discord4j.core.spec.InteractionApplicationCommandCallbackReplyMono
import discord4j.discordjson.json.ApplicationCommandRequest

interface SlashCommand {

    fun execute(event: ApplicationCommandInteractionEvent): InteractionApplicationCommandCallbackReplyMono

    fun createCommand(): ApplicationCommandRequest

}
package com.wittano.komputer.command

import discord4j.core.event.domain.interaction.ApplicationCommandInteractionEvent
import discord4j.core.`object`.command.ApplicationCommandOption
import discord4j.core.spec.InteractionApplicationCommandCallbackReplyMono
import discord4j.discordjson.json.ApplicationCommandOptionData
import discord4j.discordjson.json.ApplicationCommandRequest

class JokeCommand : SlashCommand {
    override fun execute(event: ApplicationCommandInteractionEvent): InteractionApplicationCommandCallbackReplyMono {
        TODO("Not yet implemented")
    }

    override fun createCommand(): ApplicationCommandRequest {
        return ApplicationCommandRequest.builder()
            .name("joke")
            .description("Tell me some joke")
            .options(
                listOf(
                    ApplicationCommandOptionData.builder()
                        .type(ApplicationCommandOption.Type.STRING.value)
                        .required(false)
                        .name("category")
                        .description("Joke category")
                        .choices(JOKE_CATEGORIES)
                        .build(),
                    ApplicationCommandOptionData.builder()
                        .type(ApplicationCommandOption.Type.STRING.value)
                        .required(false)
                        .name("type")
                        .description("Type of joke")
                        .choices(JOKE_TYPES)
                        .build()
                )
            )
            .build()
    }
}
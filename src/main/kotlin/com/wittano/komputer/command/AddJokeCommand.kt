package com.wittano.komputer.command

import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import discord4j.core.`object`.command.ApplicationCommandOption
import discord4j.discordjson.json.ApplicationCommandOptionData
import discord4j.discordjson.json.ApplicationCommandRequest
import reactor.core.publisher.Mono

class AddJokeCommand : SlashCommand {
    override fun execute(event: ChatInputInteractionEvent): Mono<Void> {
        TODO("Not yet implemented")
    }

    override fun createCommand(): ApplicationCommandRequest {
        return ApplicationCommandRequest.builder()
            .name("add-joke")
            .description("Add new joke to server database")
            .options(
                listOf(
                    ApplicationCommandOptionData.builder()
                        .type(ApplicationCommandOption.Type.STRING.value)
                        .required(true)
                        .name("category")
                        .description("Joke category")
                        .choices(JOKE_CATEGORIES)
                        .build(),
                    ApplicationCommandOptionData.builder()
                        .type(ApplicationCommandOption.Type.STRING.value)
                        .required(true)
                        .name("type")
                        .description("Type of joke")
                        .choices(JOKE_TYPES)
                        .build(),
                    ApplicationCommandOptionData.builder()
                        .type(ApplicationCommandOption.Type.STRING.value)
                        .required(true)
                        .name("content")
                        .description("Joke content")
                        .build(),
                    ApplicationCommandOptionData.builder()
                        .type(ApplicationCommandOption.Type.STRING.value)
                        .required(false)
                        .name("question")
                        .description("Question part of joke")
                        .build()
                )
            )
            .build()
    }

}
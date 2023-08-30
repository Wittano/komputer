package com.wittano.komputer.command

import com.google.inject.Inject
import com.wittano.komputer.joke.JokeCategory
import com.wittano.komputer.joke.JokeType
import com.wittano.komputer.joke.jokedev.JokeDevClient
import com.wittano.komputer.message.createJokeMessage
import com.wittano.komputer.message.createJokeReactionButtons
import com.wittano.komputer.utils.toNullable
import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import discord4j.core.`object`.command.ApplicationCommandInteractionOption
import discord4j.core.`object`.command.ApplicationCommandOption
import discord4j.core.`object`.component.ActionRow
import discord4j.core.spec.InteractionApplicationCommandCallbackSpec
import discord4j.discordjson.json.ApplicationCommandOptionData
import discord4j.discordjson.json.ApplicationCommandRequest
import reactor.core.publisher.Mono

class JokeCommand @Inject constructor(
    private val jokeDevClient: JokeDevClient
) : SlashCommand {
    override fun execute(event: ChatInputInteractionEvent): Mono<Void> {
        val category = event.getOption("category")
            .flatMap(ApplicationCommandInteractionOption::getValue)
            .toNullable()
            ?.asString()
            ?.let { category -> JokeCategory.entries.find { it.category == category } }
            ?: JokeCategory.ANY

        val type = event.getOption("type")
            .flatMap(ApplicationCommandInteractionOption::getValue)
            .toNullable()
            ?.asString()
            ?.let { type -> JokeType.entries.find { it.value == type } }
            ?: JokeType.SINGLE

        val joke = jokeDevClient.getRandomJoke(category, type)

        return event.reply(
            InteractionApplicationCommandCallbackSpec.builder()
                .addEmbed(createJokeMessage(joke))
                .addComponent(ActionRow.of(createJokeReactionButtons()))
                .build()
        )
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
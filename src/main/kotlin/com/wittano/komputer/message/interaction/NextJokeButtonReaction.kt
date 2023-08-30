package com.wittano.komputer.message.interaction

import com.google.inject.Inject
import com.wittano.komputer.joke.JokeCategory
import com.wittano.komputer.joke.JokeType
import com.wittano.komputer.joke.jokedev.JokeDevApiException
import com.wittano.komputer.joke.jokedev.JokeDevClient
import com.wittano.komputer.message.createErrorMessage
import com.wittano.komputer.message.createJokeMessage
import com.wittano.komputer.message.createJokeReactionButtons
import com.wittano.komputer.utils.toNullable
import discord4j.core.event.domain.interaction.ButtonInteractionEvent
import discord4j.core.`object`.component.ActionRow
import discord4j.core.spec.InteractionApplicationCommandCallbackSpec
import reactor.core.publisher.Mono

class NextJokeButtonReaction @Inject constructor(
    private val jokeDevClient: JokeDevClient
) : ButtonReaction {
    override fun execute(event: ButtonInteractionEvent): Mono<Void> {
        val fields = event.interaction.message.toNullable()?.embeds?.get(0)?.fields
        val category = fields?.get(fields.size - 1)
            ?.value
            ?.let { value -> JokeCategory.entries.find { it.polishTranslate == value } }
            ?: JokeCategory.ANY

        val type = if (fields?.size == 3) {
            JokeType.TWO_PART
        } else {
            JokeType.SINGLE
        }

        val joke = try {
            jokeDevClient.getRandomJoke(category, type)
        } catch (_: JokeDevApiException) {
            return event.reply(createErrorMessage())
        }

        return event.reply(
            InteractionApplicationCommandCallbackSpec.builder()
                .addEmbed(createJokeMessage(joke))
                .addComponent(ActionRow.of(createJokeReactionButtons()))
                .build()
        )
    }
}
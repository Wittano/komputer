package com.wittano.komputer.command

import com.wittano.komputer.joke.JokeCategory
import com.wittano.komputer.joke.JokeType
import com.wittano.komputer.joke.jokedev.JokeDevApiException
import com.wittano.komputer.joke.jokedev.JokeDevClient
import com.wittano.komputer.message.createErrorMessage
import com.wittano.komputer.message.createJokeMessage
import com.wittano.komputer.message.createJokeReactionButtons
import com.wittano.komputer.utils.getJokeCategory
import com.wittano.komputer.utils.getJokeType
import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import discord4j.core.`object`.component.ActionRow
import discord4j.core.spec.InteractionApplicationCommandCallbackSpec
import reactor.core.publisher.Mono
import javax.inject.Inject

class JokeCommand @Inject constructor(
    private val jokeDevClient: JokeDevClient
) : SlashCommand {
    override fun execute(event: ChatInputInteractionEvent): Mono<Void> {
        val category = event.getJokeCategory() ?: JokeCategory.ANY
        val type = event.getJokeType() ?: JokeType.SINGLE

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
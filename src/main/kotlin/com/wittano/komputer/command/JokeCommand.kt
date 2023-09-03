package com.wittano.komputer.command

import com.wittano.komputer.joke.JokeApiService
import com.wittano.komputer.joke.JokeCategory
import com.wittano.komputer.joke.JokeRandomService
import com.wittano.komputer.joke.JokeType
import com.wittano.komputer.joke.api.jokedev.JokeDevApiException
import com.wittano.komputer.message.createJokeMessage
import com.wittano.komputer.message.createJokeReactionButtons
import com.wittano.komputer.message.resource.ErrorMessage
import com.wittano.komputer.utils.getJokeCategory
import com.wittano.komputer.utils.getJokeType
import com.wittano.komputer.utils.getRandomJoke
import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import discord4j.core.`object`.component.ActionRow
import discord4j.core.spec.InteractionApplicationCommandCallbackSpec
import reactor.core.publisher.Mono
import reactor.core.scheduler.Schedulers
import javax.inject.Inject

class JokeCommand @Inject constructor(
    private val jokeDevClient: JokeApiService,
    private val jokeRandomServices: Set<@JvmSuppressWildcards JokeRandomService>
) : SlashCommand {
    override fun execute(event: ChatInputInteractionEvent): Mono<Void> {
        val category = event.getJokeCategory() ?: JokeCategory.ANY
        val type = event.getJokeType() ?: JokeType.SINGLE

        if (!jokeDevClient.supports(type)) {
            return Mono.error(JokeDevApiException("Joke type '$type' isn't support", ErrorMessage.UNSUPPORTED_TYPE))
        }

        return getRandomJoke(type, category, jokeRandomServices)
            .publishOn(Schedulers.boundedElastic())
            .flatMap {
                val message = InteractionApplicationCommandCallbackSpec.builder()
                    .addEmbed(createJokeMessage(it))
                    .addComponent(ActionRow.of(createJokeReactionButtons()))
                    .build()

                event.reply(message)
            }

    }
}
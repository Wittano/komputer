package com.wittano.komputer.core.message.interaction

import com.wittano.komputer.core.joke.JokeApiService
import com.wittano.komputer.core.joke.JokeCategory
import com.wittano.komputer.core.joke.JokeRandomService
import com.wittano.komputer.core.joke.JokeType
import com.wittano.komputer.core.joke.api.jokedev.JokeDevApiException
import com.wittano.komputer.core.message.createJokeMessage
import com.wittano.komputer.core.message.createJokeReactionButtons
import com.wittano.komputer.core.message.resource.ButtonLabel
import com.wittano.komputer.core.message.resource.ErrorMessage
import com.wittano.komputer.core.message.resource.getButtonLabel
import com.wittano.komputer.core.message.resource.toLocale
import com.wittano.komputer.core.utils.getRandomJoke
import discord4j.core.event.domain.interaction.ButtonInteractionEvent
import discord4j.core.`object`.component.ActionRow
import discord4j.core.spec.InteractionApplicationCommandCallbackSpec
import reactor.core.publisher.Mono
import javax.inject.Inject
import kotlin.jvm.optionals.getOrNull
import kotlin.random.Random

class NextJokeButtonReaction @Inject constructor(
    private val jokeDevClient: JokeApiService,
    private val jokeRandomServices: Set<@JvmSuppressWildcards JokeRandomService>
) : ButtonReaction {
    override fun execute(event: ButtonInteractionEvent): Mono<Void> {
        val fields = event.interaction.message.getOrNull()?.embeds?.get(0)?.fields
        val category = fields?.get(fields.size - 1)
            ?.value
            ?.let { value -> JokeCategory.entries.find { it.polishTranslate == value } }
            ?.takeIf { it != JokeCategory.YO_MAMA }
            ?: JokeCategory.ANY

        val type = if (fields?.size == 3) {
            JokeType.TWO_PART
        } else {
            JokeType.SINGLE
        }

        if (!jokeDevClient.supports(type)) {
            return Mono.error(
                JokeDevApiException(
                    "Joke type '$type' isn't support by API",
                    ErrorMessage.UNSUPPORTED_TYPE
                )
            )
        }

        val apologies = getButtonLabel(ButtonLabel.APOLOGIES, event.interaction.userLocale.toLocale())
            .takeIf {
                Random.nextInt().mod(7) == 0
            } ?: ""

        return getRandomJoke(type, category, jokeRandomServices).flatMap {
            event.reply(
                InteractionApplicationCommandCallbackSpec.builder()
                    .content(apologies)
                    .addEmbed(createJokeMessage(it))
                    .addComponent(ActionRow.of(createJokeReactionButtons(event.interaction.userLocale.toLocale())))
                    .build()
            )
        }
    }
}
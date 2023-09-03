package com.wittano.komputer.message.interaction

import com.wittano.komputer.joke.JokeApiService
import com.wittano.komputer.joke.JokeCategory
import com.wittano.komputer.joke.JokeRandomService
import com.wittano.komputer.joke.JokeType
import com.wittano.komputer.joke.api.jokedev.JokeDevApiException
import com.wittano.komputer.message.createErrorMessage
import com.wittano.komputer.message.createJokeMessage
import com.wittano.komputer.message.createJokeReactionButtons
import com.wittano.komputer.message.resource.ErrorMessage
import com.wittano.komputer.utils.filterService
import discord4j.core.event.domain.interaction.ButtonInteractionEvent
import discord4j.core.`object`.component.ActionRow
import discord4j.core.spec.InteractionApplicationCommandCallbackSpec
import org.slf4j.LoggerFactory
import reactor.core.publisher.Mono
import reactor.core.scheduler.Schedulers
import javax.inject.Inject
import kotlin.jvm.optionals.getOrNull
import kotlin.random.Random

class NextJokeButtonReaction @Inject constructor(
    private val jokeDevClient: JokeApiService,
    private val jokeRandomService: Set<@JvmSuppressWildcards JokeRandomService>
) : ButtonReaction {
    private val log = LoggerFactory.getLogger(this::class.qualifiedName)

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

        val joke = jokeRandomService.filterService(type, category).random().getRandom(category, type)
        val apologies = "Przepraszam".takeIf { Random.nextInt().mod(7) == 0 } ?: ""

        return joke.flatMap {
            event.reply(
                InteractionApplicationCommandCallbackSpec.builder()
                    .content(apologies)
                    .addEmbed(createJokeMessage(it))
                    .addComponent(ActionRow.of(createJokeReactionButtons()))
                    .build()
            )
        }.publishOn(Schedulers.boundedElastic())
            .onErrorResume {
                log.error("Unexpected error during response on next joke reaction button", it)

                event.reply(createErrorMessage())
            }
    }
}
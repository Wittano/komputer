package com.wittano.komputer.bot.message.interaction

import com.wittano.komputer.bot.joke.JokeCategory
import com.wittano.komputer.bot.joke.JokeRandomService
import com.wittano.komputer.bot.message.createJokeMessage
import com.wittano.komputer.bot.message.createJokeReactionButtons
import com.wittano.komputer.bot.utils.getRandomJoke
import com.wittano.komputer.bot.utils.joke.getGuid
import com.wittano.komputer.bot.utils.mongodb.getGlobalLanguage
import com.wittano.komputer.commons.transtation.ButtonLabel
import com.wittano.komputer.commons.transtation.getButtonLabel
import discord4j.core.event.domain.interaction.ButtonInteractionEvent
import discord4j.core.`object`.component.ActionRow
import discord4j.core.spec.InteractionApplicationCommandCallbackSpec
import reactor.core.publisher.Mono
import javax.inject.Inject
import kotlin.jvm.optionals.getOrNull
import kotlin.random.Random

class NextJokeButtonReaction @Inject constructor(
    private val jokeRandomServices: Set<@JvmSuppressWildcards JokeRandomService>
) : ButtonReaction {
    override fun execute(event: ButtonInteractionEvent): Mono<Void> {
        val fields = event.interaction.message.getOrNull()?.embeds?.get(0)?.fields
        val category = fields?.get(fields.size - 1)
            ?.value
            ?.let { value -> JokeCategory.entries.find { it.polishTranslate == value } }
            ?.takeIf { it != JokeCategory.YO_MAMA }

        val language = getGlobalLanguage(event.getGuid())
        val apologies = getButtonLabel(ButtonLabel.APOLOGIES, language)
            .takeIf { Random.nextInt() % 7 == 0 }
            .orEmpty()

        return getRandomJoke(
            null,
            category,
            jokeRandomServices,
            null
        ).flatMap {
            event.reply(
                InteractionApplicationCommandCallbackSpec.builder()
                    .content(apologies)
                    .addEmbed(createJokeMessage(it, getGlobalLanguage(event.getGuid())))
                    .addComponent(ActionRow.of(createJokeReactionButtons(language)))
                    .build()
            )
        }
    }
}
package com.wittano.komputer.bot.message.interaction

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
import kotlin.random.Random

class NextRandomJokeButtonReaction @Inject constructor(
    private val jokeRandomServices: Set<@JvmSuppressWildcards JokeRandomService>
) : ButtonReaction {
    override fun execute(event: ButtonInteractionEvent): Mono<Void> {
        val language = getGlobalLanguage(event.getGuid())
        val apologies = getButtonLabel(ButtonLabel.APOLOGIES, language)
            .takeIf { Random.nextInt() % 7 == 0 }
            .orEmpty()

        return getRandomJoke(
            null,
            null,
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
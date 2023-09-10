package com.wittano.komputer.bot.command

import com.wittano.komputer.bot.joke.JokeApiService
import com.wittano.komputer.bot.joke.JokeRandomService
import com.wittano.komputer.bot.joke.api.jokedev.JokeDevApiException
import com.wittano.komputer.bot.message.createJokeMessage
import com.wittano.komputer.bot.message.createJokeReactionButtons
import com.wittano.komputer.bot.utils.getRandomJoke
import com.wittano.komputer.bot.utils.joke.getGuid
import com.wittano.komputer.bot.utils.joke.getJokeCategory
import com.wittano.komputer.bot.utils.joke.getJokeType
import com.wittano.komputer.bot.utils.joke.getLanguageOptional
import com.wittano.komputer.bot.utils.mongodb.getGlobalLanguage
import com.wittano.komputer.commons.transtation.ErrorMessage
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
        val category = event.getJokeCategory()
        val type = event.getJokeType()
        val language = event.getLanguageOptional()

        if (type != null && !jokeDevClient.supports(type)) {
            return Mono.error(JokeDevApiException("Joke type '$type' isn't support", ErrorMessage.UNSUPPORTED_TYPE))
        }

        return getRandomJoke(type, category, jokeRandomServices, language)
            .publishOn(Schedulers.boundedElastic())
            .flatMap {
                val message = InteractionApplicationCommandCallbackSpec.builder()
                    .addEmbed(createJokeMessage(it, getGlobalLanguage(event.getGuid())))
                    .addComponent(ActionRow.of(createJokeReactionButtons(getGlobalLanguage(event.getGuid()))))
                    .build()

                event.reply(message)
            }

    }
}
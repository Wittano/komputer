package com.wittano.komputer.bot.command

import com.wittano.komputer.bot.joke.Joke
import com.wittano.komputer.bot.joke.JokeType
import com.wittano.komputer.bot.joke.mongodb.JokeDatabaseService
import com.wittano.komputer.bot.message.createErrorMessage
import com.wittano.komputer.bot.utils.getJokeCategory
import com.wittano.komputer.bot.utils.getJokeType
import com.wittano.komputer.commons.transtation.SuccessfulMessage
import com.wittano.komputer.commons.transtation.getSuccessfulMessage
import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import discord4j.core.`object`.command.ApplicationCommandInteractionOption
import discord4j.core.spec.InteractionApplicationCommandCallbackSpec
import org.slf4j.LoggerFactory
import reactor.core.publisher.Mono
import java.time.Duration
import javax.inject.Inject
import kotlin.jvm.optionals.getOrNull

class AddJokeCommand @Inject constructor(
    private val databaseService: JokeDatabaseService
) : SlashCommand {
    private val log = LoggerFactory.getLogger(this::class.qualifiedName)

    override fun execute(event: ChatInputInteractionEvent): Mono<Void> {
        val joke = try {
            val content = event.getOption("content")
                .flatMap(ApplicationCommandInteractionOption::getValue)
                .filter { it.asString().isNotBlank() }
                .getOrNull()
                ?.asString()

            val question = event.getOption("question")
                .flatMap(ApplicationCommandInteractionOption::getValue)
                .getOrNull()
                ?.asString()

            val jokeType = event.getJokeType() ?: JokeType.SINGLE

            Joke(
                category = event.getJokeCategory()!!,
                type = jokeType,
                answer = content!!,
                question = question
            )
        } catch (_: NullPointerException) {
            log.warn("During getting joke data throw unexpected error. Some required field is missing")
            return event.reply(createErrorMessage())
        }

        return databaseService.add(joke)
            .timeout(Duration.ofSeconds(1))
            .flatMap { sendPositiveFeedback(it, event) }
    }

    private fun sendPositiveFeedback(jokeId: String, event: ChatInputInteractionEvent): Mono<Void> {
        val messageResponse = InteractionApplicationCommandCallbackSpec.builder()
            .content(getSuccessfulMessage(SuccessfulMessage.ADD_JOKE).format(jokeId))
            .build()
            .withEphemeral(true)

        return event.reply(messageResponse)
    }
}
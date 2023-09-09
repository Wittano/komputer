package com.wittano.komputer.bot.command

import com.wittano.komputer.bot.joke.Joke
import com.wittano.komputer.bot.joke.JokeCategory
import com.wittano.komputer.bot.joke.JokeException
import com.wittano.komputer.bot.joke.JokeType
import com.wittano.komputer.bot.joke.mongodb.JokeDatabaseService
import com.wittano.komputer.bot.utils.getJokeCategory
import com.wittano.komputer.bot.utils.getJokeType
import com.wittano.komputer.bot.utils.getLanguage
import com.wittano.komputer.commons.transtation.ErrorMessage
import com.wittano.komputer.commons.transtation.SuccessfulMessage
import com.wittano.komputer.commons.transtation.getSuccessfulMessage
import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import discord4j.core.`object`.command.ApplicationCommandInteractionOption
import discord4j.core.spec.InteractionApplicationCommandCallbackSpec
import reactor.core.publisher.Mono
import java.time.Duration
import javax.inject.Inject
import kotlin.jvm.optionals.getOrNull

class AddJokeCommand @Inject constructor(
    private val databaseService: JokeDatabaseService
) : SlashCommand {
    override fun execute(event: ChatInputInteractionEvent): Mono<Void> {
        val content = event.getOption("content")
            .flatMap(ApplicationCommandInteractionOption::getValue)
            .filter { it.asString().isNotBlank() }
            .getOrNull()
            ?.asString()
            .orEmpty()

        val jokeType = event.getJokeType() ?: JokeType.SINGLE
        val question = event.getOption("question")
            .flatMap(ApplicationCommandInteractionOption::getValue)
            .getOrNull()
            ?.asString()
            ?.takeIf { jokeType != JokeType.SINGLE }
            ?: return Mono.error(
                JokeException(
                    "Question part in Two-Part joke is required!",
                    ErrorMessage.MISSING_QUESTION_FILED
                )
            )

        val joke = Joke(
            category = event.getJokeCategory() ?: JokeCategory.ANY,
            type = jokeType,
            answer = content,
            question = question,
            language = event.getLanguage()
        )

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
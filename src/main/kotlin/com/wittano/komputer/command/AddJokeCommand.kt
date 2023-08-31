package com.wittano.komputer.command

import com.wittano.komputer.joke.Joke
import com.wittano.komputer.joke.JokeType
import com.wittano.komputer.joke.mongodb.JokeDatabaseManager
import com.wittano.komputer.message.createErrorMessage
import com.wittano.komputer.utils.getJokeCategory
import com.wittano.komputer.utils.getJokeType
import com.wittano.komputer.utils.toNullable
import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import discord4j.core.`object`.command.ApplicationCommandInteractionOption
import discord4j.core.spec.InteractionApplicationCommandCallbackSpec
import org.slf4j.LoggerFactory
import reactor.core.publisher.Mono
import java.time.Duration
import javax.inject.Inject

class AddJokeCommand @Inject constructor(
    private val databaseService: JokeDatabaseManager
) : SlashCommand {
    private val log = LoggerFactory.getLogger(this::class.qualifiedName)

    override fun execute(event: ChatInputInteractionEvent): Mono<Void> {
        val joke = try {
            val content = event.getOption("content")
                .flatMap(ApplicationCommandInteractionOption::getValue)
                .filter {
                    it.asString().isNotBlank()
                }
                .toNullable()
                ?.asString()

            val question = event.getOption("question")
                .flatMap(ApplicationCommandInteractionOption::getValue)
                .toNullable()
                ?.asString()

            val jokeType = event.getJokeType()

            Joke(
                category = event.getJokeCategory()!!,
                type = jokeType!!,
                answer = content!!,
                question = if (jokeType == JokeType.TWO_PART) {
                    question!!
                } else {
                    null
                }
            )
        } catch (_: NullPointerException) {
            log.warn("During getting joke data throw unexpected error. Some required field is missing")
            return event.reply(createErrorMessage())
        }

        return databaseService.add(joke)
            .timeout(Duration.ofSeconds(1))
            .flatMap { sendPositiveFeedback(it, event) }
            .switchIfEmpty(event.reply("BEEP BOOP. Coś poszło nie tak"))
    }

    private fun sendPositiveFeedback(jokeId: String, event: ChatInputInteractionEvent): Mono<Void> {
        val messageResponse = InteractionApplicationCommandCallbackSpec.builder()
            .content("BEEP BOOP. Udało się dodać żart. Twój żart ma id: $jokeId")
            .build()
            .withEphemeral(true)

        return event.reply(messageResponse)
    }
}
package com.wittano.komputer.core.command

import com.wittano.komputer.core.joke.Joke
import com.wittano.komputer.core.joke.JokeType
import com.wittano.komputer.core.message.createErrorMessage
import com.wittano.komputer.core.utils.getJokeCategory
import com.wittano.komputer.core.utils.getJokeType
import com.wittano.komputer.joke.mongodb.JokeDatabaseService
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
                .filter {
                    it.asString().isNotBlank()
                }.getOrNull()
                ?.asString()

            val question = event.getOption("question")
                .flatMap(ApplicationCommandInteractionOption::getValue)
                .getOrNull()
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
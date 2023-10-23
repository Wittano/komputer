package com.wittano.komputer.bot.command

import com.wittano.komputer.bot.joke.mongodb.JokeDatabaseService
import com.wittano.komputer.commons.transtation.SuccessfulMessage
import com.wittano.komputer.commons.transtation.getSuccessfulMessage
import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import discord4j.core.`object`.command.ApplicationCommandInteractionOption
import reactor.core.publisher.Mono
import java.time.Duration
import javax.inject.Inject

class RemoveJokeCommand @Inject constructor(
    private val jokeDatabaseService: JokeDatabaseService
) : SlashCommand {
    override fun execute(event: ChatInputInteractionEvent): Mono<Void> {
        val jokeId = event.getOption("id")
            .flatMap(ApplicationCommandInteractionOption::getValue)
            .get()
            .asString()

        return jokeDatabaseService.remove(jokeId)
            .timeout(Duration.ofSeconds(1))
            .then(event.reply(getSuccessfulMessage(SuccessfulMessage.REMOVE_JOKE)))
    }
}
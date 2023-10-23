package com.wittano.komputer.bot.joke

import reactor.core.publisher.Mono
import java.util.*

fun interface JokeRandomService {
    fun getRandom(
        category: JokeCategory?,
        type: JokeType?,
        language: Locale?
    ): Mono<Joke>
}
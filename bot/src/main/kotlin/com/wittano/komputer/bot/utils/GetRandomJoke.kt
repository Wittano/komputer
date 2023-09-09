package com.wittano.komputer.bot.utils

import com.wittano.komputer.bot.joke.*
import com.wittano.komputer.bot.joke.api.rapidapi.RapidApiException
import com.wittano.komputer.commons.transtation.ErrorMessage
import reactor.core.publisher.Mono

internal fun getRandomJoke(
    type: JokeType,
    category: JokeCategory,
    jokeRandomServices: Set<JokeRandomService>
): Mono<Joke> {
    val jokeRandomService = jokeRandomServices.filterService(type, category)
    return jokeRandomService.random().getRandom(category, type)
        .onErrorResume {
            if (it is RapidApiException) {
                return@onErrorResume jokeRandomService.excludeRapidApiServices()
                    .takeIf { list -> list.isNotEmpty() }
                    ?.random()
                    ?.getRandom(category, type)
                    ?: Mono.error(JokeException("Joke not found", ErrorMessage.JOKE_NOT_FOUND))
            }

            Mono.error(it)
        }
}
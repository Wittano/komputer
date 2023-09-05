package com.wittano.komputer.core.utils

import com.wittano.komputer.core.joke.*
import com.wittano.komputer.core.joke.api.rapidapi.RapidApiException
import com.wittano.komputer.core.message.resource.ErrorMessage
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
package com.wittano.komputer.bot.utils

import com.wittano.komputer.bot.joke.*
import com.wittano.komputer.bot.joke.api.rapidapi.RapidApiException
import com.wittano.komputer.commons.transtation.ErrorMessage
import reactor.core.publisher.Mono
import java.util.*

internal fun getRandomJoke(
    type: JokeType,
    category: JokeCategory,
    jokeRandomServices: Set<JokeRandomService>,
    language: Locale = Locale.ENGLISH
): Mono<Joke> {
    val jokeNotFoundError = Mono.error<Joke>(JokeException("Joke not found", ErrorMessage.JOKE_NOT_FOUND))
    val jokeRandomService = jokeRandomServices.filterService(type, category, language)

    return jokeRandomService.takeIf { it.isNotEmpty() }
        ?.random()
        ?.getRandom(category, type, language)
        ?.onErrorResume {
            if (it is RapidApiException) {
                return@onErrorResume jokeRandomService.excludeRapidApiServices()
                    .takeIf { list -> list.isNotEmpty() }
                    ?.random()
                    ?.getRandom(category, type, language)
                    ?: jokeNotFoundError
            }

            Mono.error(it)
        }
        ?: jokeNotFoundError
}
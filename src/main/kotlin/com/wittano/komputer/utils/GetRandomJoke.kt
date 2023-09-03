package com.wittano.komputer.utils

import com.wittano.komputer.joke.Joke
import com.wittano.komputer.joke.JokeCategory
import com.wittano.komputer.joke.JokeRandomService
import com.wittano.komputer.joke.JokeType
import com.wittano.komputer.joke.api.rapidapi.RapidApiException
import reactor.core.publisher.Mono

fun getRandomJoke(
    type: JokeType,
    category: JokeCategory,
    jokeRandomServices: Set<JokeRandomService>
): Mono<Joke> {
    val jokeRandomService = jokeRandomServices.filterService(type, category)
    return jokeRandomService.random().getRandom(category, type)
        .onErrorResume {
            if (it is RapidApiException) {
                return@onErrorResume jokeRandomService.excludeRapidApiServices()
                    .random()
                    .getRandom(category, type)
            }

            Mono.error(it)
        }
}
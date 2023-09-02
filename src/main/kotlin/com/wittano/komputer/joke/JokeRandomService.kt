package com.wittano.komputer.joke

import reactor.core.publisher.Mono

fun interface JokeRandomService {
    fun getRandom(category: JokeCategory?, type: JokeType): Mono<Joke>
}
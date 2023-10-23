package com.wittano.komputer.bot.joke

import reactor.core.publisher.Mono

interface JokeService {

    fun add(joke: Joke): Mono<String>

    fun remove(id: String): Mono<Void>

    fun get(id: String): Mono<Joke>

}
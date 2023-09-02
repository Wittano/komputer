package com.wittano.komputer.joke

import reactor.core.publisher.Mono

// TODO Add more external API intergration
interface JokeApiService {

    fun getRandom(category: JokeCategory, type: JokeType): Mono<Joke>
    fun supports(type: JokeType): Boolean
    fun supports(category: JokeCategory): Boolean

}
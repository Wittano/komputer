package com.wittano.komputer.joke

// TODO Add more external API intergration
interface JokeApiService {

    fun supports(type: JokeType): Boolean = true
    fun supports(category: JokeCategory): Boolean

}
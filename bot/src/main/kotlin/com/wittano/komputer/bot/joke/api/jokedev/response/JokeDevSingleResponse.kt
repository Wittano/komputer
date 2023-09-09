package com.wittano.komputer.bot.joke.api.jokedev.response

import com.wittano.komputer.bot.joke.Joke
import com.wittano.komputer.bot.joke.JokeExtractor
import com.wittano.komputer.bot.joke.JokeType

data class JokeDevSingleResponse(
    val category: JokeDevCategory,
    val error: Boolean,
    val flags: Flags,
    val id: Int,
    val joke: String,
    val lang: String,
    val safe: Boolean,
    val type: String
) : JokeExtractor {
    override fun toJoke(): Joke = Joke(
        category = category.jokeCategory,
        answer = joke,
        type = JokeType.SINGLE
    )
}
package com.wittano.komputer.joke.jokedev.response

import com.wittano.komputer.joke.Joke
import com.wittano.komputer.joke.JokeExtractor
import com.wittano.komputer.joke.JokeType
import com.wittano.komputer.joke.jokedev.toJokeCategory

data class JokeDevTwoPartResponse(
    val category: String,
    val delivery: String,
    val error: Boolean,
    val flags: Flags,
    val id: Int,
    val lang: String,
    val safe: Boolean,
    val setup: String,
    val type: String
) : JokeExtractor {
    override fun toJoke(): Joke = Joke(
        content = delivery,
        question = setup,
        category = category.toJokeCategory(),
        type = JokeType.TWO_PART
    )
}
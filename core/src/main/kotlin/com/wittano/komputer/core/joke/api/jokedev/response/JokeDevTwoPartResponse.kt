package com.wittano.komputer.core.joke.api.jokedev.response

import com.wittano.komputer.core.joke.Joke
import com.wittano.komputer.core.joke.JokeExtractor
import com.wittano.komputer.core.joke.JokeType
import com.wittano.komputer.core.joke.api.jokedev.toJokeCategory

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
        answer = delivery,
        question = setup,
        category = category.toJokeCategory(),
        type = JokeType.TWO_PART
    )
}
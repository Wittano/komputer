package com.wittano.komputer.bot.joke.api.jokedev

import com.wittano.komputer.bot.joke.JokeCategory

internal fun String.toJokeCategory(): JokeCategory = try {
    JokeCategory.valueOf(this)
} catch (ex: IllegalArgumentException) {
    JokeCategory.ANY
}
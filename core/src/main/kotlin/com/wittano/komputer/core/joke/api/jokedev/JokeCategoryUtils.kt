package com.wittano.komputer.core.joke.api.jokedev

import com.wittano.komputer.core.joke.JokeCategory

internal fun String.toJokeCategory(): JokeCategory = try {
    JokeCategory.valueOf(this)
} catch (ex: IllegalArgumentException) {
    JokeCategory.ANY
}
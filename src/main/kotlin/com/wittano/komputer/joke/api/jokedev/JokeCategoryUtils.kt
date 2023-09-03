package com.wittano.komputer.joke.api.jokedev

import com.wittano.komputer.joke.JokeCategory

internal fun String.toJokeCategory(): JokeCategory = try {
    JokeCategory.valueOf(this)
} catch (ex: IllegalArgumentException) {
    JokeCategory.ANY
}
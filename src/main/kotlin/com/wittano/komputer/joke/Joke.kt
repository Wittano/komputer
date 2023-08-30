package com.wittano.komputer.joke

data class Joke(
    val content: String,
    val category: JokeCategory,
    val type: JokeType,
    val question: String? = null
)

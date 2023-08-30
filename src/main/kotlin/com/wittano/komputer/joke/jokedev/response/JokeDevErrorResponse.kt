package com.wittano.komputer.joke.jokedev.response

data class JokeDevErrorResponse(
    val additionalInfo: String,
    val causedBy: List<String>,
    val code: Int,
    val error: Boolean,
    val internalError: Boolean,
    val message: String,
    val timestamp: Long
)
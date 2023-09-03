package com.wittano.komputer.core.joke.api.jokedev

import com.wittano.komputer.core.joke.api.JokeApiException
import com.wittano.komputer.core.joke.api.jokedev.response.JokeDevErrorResponse
import com.wittano.komputer.core.message.resource.ErrorMessage

class JokeDevApiException(message: String, code: ErrorMessage, val response: JokeDevErrorResponse? = null) :
    JokeApiException(message, code)

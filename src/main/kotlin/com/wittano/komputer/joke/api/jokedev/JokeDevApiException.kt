package com.wittano.komputer.joke.api.jokedev

import com.wittano.komputer.joke.api.JokeApiException
import com.wittano.komputer.joke.api.jokedev.response.JokeDevErrorResponse
import com.wittano.komputer.message.resource.ErrorMessage

class JokeDevApiException(message: String, code: ErrorMessage, val response: JokeDevErrorResponse? = null) :
    JokeApiException(message, code)

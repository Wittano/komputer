package com.wittano.komputer.joke.jokedev

import com.wittano.komputer.joke.JokeApiException
import com.wittano.komputer.joke.jokedev.response.JokeDevErrorResponse
import com.wittano.komputer.message.resource.ErrorMessage

class JokeDevApiException(message: String, code: ErrorMessage, val response: JokeDevErrorResponse? = null) :
    JokeApiException(message, code)

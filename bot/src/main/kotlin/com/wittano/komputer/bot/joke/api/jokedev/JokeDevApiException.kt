package com.wittano.komputer.bot.joke.api.jokedev

import com.wittano.komputer.bot.joke.api.JokeApiException
import com.wittano.komputer.bot.joke.api.jokedev.response.JokeDevErrorResponse
import com.wittano.komputer.commons.transtation.ErrorMessage

class JokeDevApiException(message: String, code: ErrorMessage, val response: JokeDevErrorResponse? = null) :
    JokeApiException(message, code)

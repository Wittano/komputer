package com.wittano.komputer.joke.exception

import com.wittano.komputer.message.resource.ErrorMessage

open class JokeApiException(msg: String, val code: ErrorMessage) : RuntimeException(msg)
package com.wittano.komputer.joke.api

import com.wittano.komputer.joke.JokeException
import com.wittano.komputer.message.resource.ErrorMessage

open class JokeApiException(
    msg: String,
    code: ErrorMessage,
) : JokeException(msg, code)
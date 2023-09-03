package com.wittano.komputer.core.joke.api

import com.wittano.komputer.core.joke.JokeException
import com.wittano.komputer.core.message.resource.ErrorMessage

open class JokeApiException(
    msg: String,
    code: ErrorMessage,
) : JokeException(msg, code)
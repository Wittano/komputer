package com.wittano.komputer.bot.joke.api

import com.wittano.komputer.bot.joke.JokeException
import com.wittano.komputer.commons.transtation.ErrorMessage

open class JokeApiException(
    msg: String,
    code: ErrorMessage,
) : JokeException(msg, code)
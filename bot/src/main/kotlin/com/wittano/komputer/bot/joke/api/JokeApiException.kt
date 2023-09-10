package com.wittano.komputer.bot.joke.api

import com.wittano.komputer.bot.joke.CommandException
import com.wittano.komputer.commons.transtation.ErrorMessage

open class JokeApiException(
    msg: String,
    code: ErrorMessage,
) : CommandException(msg, code)
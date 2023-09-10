package com.wittano.komputer.bot.joke.api

import com.wittano.komputer.bot.command.exception.CommandException
import com.wittano.komputer.commons.transtation.ErrorMessage

open class JokeApiException(
    msg: String,
    code: ErrorMessage,
) : CommandException(msg, code)
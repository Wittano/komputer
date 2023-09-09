package com.wittano.komputer.bot.joke

import com.wittano.komputer.commons.transtation.ErrorMessage

open class JokeException(msg: String, val code: ErrorMessage) : RuntimeException(msg)
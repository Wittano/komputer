package com.wittano.komputer.core.joke

import com.wittano.komputer.core.message.resource.ErrorMessage

open class JokeException(msg: String, val code: ErrorMessage) : RuntimeException(msg)
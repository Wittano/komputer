package com.wittano.komputer.joke

import com.wittano.komputer.message.resource.ErrorMessage

open class JokeException(msg: String, val code: ErrorMessage) : RuntimeException(msg)
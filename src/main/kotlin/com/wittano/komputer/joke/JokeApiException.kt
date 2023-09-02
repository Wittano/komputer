package com.wittano.komputer.joke

import com.wittano.komputer.message.resource.ErrorMessage

open class JokeApiException(msg: String, code: ErrorMessage) : JokeException(msg, code)
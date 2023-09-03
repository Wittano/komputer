package com.wittano.komputer.core.joke.mongodb

import com.wittano.komputer.core.joke.JokeException
import com.wittano.komputer.core.message.resource.ErrorMessage

class InvalidJokeIdException(msg: String, code: ErrorMessage) : JokeException(msg, code)

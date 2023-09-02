package com.wittano.komputer.joke.mongodb

import com.wittano.komputer.joke.JokeException
import com.wittano.komputer.message.resource.ErrorMessage

class InvalidJokeIdException(msg: String, code: ErrorMessage) : JokeException(msg, code)

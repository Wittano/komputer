package com.wittano.komputer.bot.joke.mongodb

import com.wittano.komputer.bot.joke.JokeException
import com.wittano.komputer.commons.transtation.ErrorMessage

class InvalidJokeIdException(msg: String, code: ErrorMessage) : JokeException(msg, code)

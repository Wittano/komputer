package com.wittano.komputer.bot.joke.mongodb

import com.wittano.komputer.bot.joke.CommandException
import com.wittano.komputer.commons.transtation.ErrorMessage

class InvalidJokeIdException(msg: String, code: ErrorMessage) : CommandException(msg, code)

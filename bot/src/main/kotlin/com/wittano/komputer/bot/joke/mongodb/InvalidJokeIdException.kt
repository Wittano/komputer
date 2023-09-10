package com.wittano.komputer.bot.joke.mongodb

import com.wittano.komputer.bot.command.exception.CommandException
import com.wittano.komputer.commons.transtation.ErrorMessage

class InvalidJokeIdException(msg: String, code: ErrorMessage) : CommandException(msg, code)

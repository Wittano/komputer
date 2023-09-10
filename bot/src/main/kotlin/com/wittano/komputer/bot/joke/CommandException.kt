package com.wittano.komputer.bot.joke

import com.wittano.komputer.commons.transtation.ErrorMessage

open class CommandException(msg: String, val code: ErrorMessage) : RuntimeException(msg)
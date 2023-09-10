package com.wittano.komputer.bot.command.exception

import com.wittano.komputer.commons.transtation.ErrorMessage

open class CommandException(
    msg: String,
    val code: ErrorMessage,
    val isUserOnlyVisible: Boolean = false
) : RuntimeException(msg)
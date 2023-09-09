package com.wittano.komputer.bot.command.exception

class CommandException(msg: String, val commandId: String) : RuntimeException(msg)
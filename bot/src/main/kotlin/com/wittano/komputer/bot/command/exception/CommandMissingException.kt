package com.wittano.komputer.bot.command.exception

class CommandMissingException(msg: String, val commandId: String) : RuntimeException(msg)
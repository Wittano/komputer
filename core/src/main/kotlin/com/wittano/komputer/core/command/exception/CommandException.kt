package com.wittano.komputer.core.command.exception

class CommandException(msg: String, val commandId: String) : RuntimeException(msg)
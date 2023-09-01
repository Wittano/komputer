package com.wittano.komputer.command.exception

class CommandException(msg: String, val commandId: String) : RuntimeException(msg)
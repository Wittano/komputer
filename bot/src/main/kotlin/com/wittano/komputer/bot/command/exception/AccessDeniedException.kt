package com.wittano.komputer.bot.command.exception

import com.wittano.komputer.commons.transtation.ErrorMessage

class AccessDeniedException(userId: String, guildId: String) : CommandException(
    "User '${userId}' hasn't enough permission to update config on '$guildId' server",
    ErrorMessage.ACCESS_DENIED,
    isUserOnlyVisible = true
)
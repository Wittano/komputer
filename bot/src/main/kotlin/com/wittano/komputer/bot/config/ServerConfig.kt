package com.wittano.komputer.bot.config

import com.wittano.komputer.commons.extensions.POLISH_LOCALE
import java.util.*

data class ServerConfig(
    val language: Locale = POLISH_LOCALE,
    val roleId: String? = null
) {
    fun toModel(guid: String) = ServerConfigModel(guid, language.language, roleId)
}

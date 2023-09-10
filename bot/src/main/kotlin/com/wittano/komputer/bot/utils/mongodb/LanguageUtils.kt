package com.wittano.komputer.bot.utils.mongodb

import com.wittano.komputer.bot.config.ConfigDatabaseService
import com.wittano.komputer.commons.extensions.POLISH_LOCALE
import reactor.core.publisher.Mono
import java.util.*
import java.util.concurrent.ConcurrentHashMap

private val languageServers = ConcurrentHashMap<String, Locale>()

internal fun getGlobalLanguage(guid: String): Locale {
    languageServers[guid]?.also { return it }

    return ConfigDatabaseService()[guid]
        .map {
            it.language
        }
        .switchIfEmpty(Mono.just(POLISH_LOCALE))
        .doOnSuccess {
            languageServers[guid] = it
        }
        .block()
        ?: POLISH_LOCALE
}
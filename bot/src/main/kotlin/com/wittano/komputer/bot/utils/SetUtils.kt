package com.wittano.komputer.bot.utils

import com.wittano.komputer.bot.joke.JokeApiService
import com.wittano.komputer.bot.joke.JokeCategory
import com.wittano.komputer.bot.joke.JokeRandomService
import com.wittano.komputer.bot.joke.JokeType
import com.wittano.komputer.bot.joke.api.rapidapi.RapidAPIService
import com.wittano.komputer.bot.joke.mongodb.JokeDatabaseService
import com.wittano.komputer.bot.utils.mongodb.isDatabaseReady
import java.util.*

internal fun Set<JokeRandomService>.filterService(
    type: JokeType?,
    category: JokeCategory?,
    language: Locale?
): List<JokeRandomService> =
    this.filter { service ->
        if (service is JokeApiService) {
            val isTypeSupport = type?.let { service.supports(it) } ?: true
            val isCategorySupport = category?.let { service.supports(it) } ?: true
            val isLanguageSupport = language?.let { service.supports(it) } ?: true

            return@filter isTypeSupport && isCategorySupport && isLanguageSupport
        }

        true
    }.filter { excludeRapidApiService(it) }
        .filter {
            if (it is JokeDatabaseService) {
                return@filter isDatabaseReady.get()
            }

            return@filter true
        }

private fun excludeRapidApiService(it: JokeRandomService): Boolean {
    if (it is RapidAPIService) {
        return it.isEnable()
    }

    return true
}

internal fun List<JokeRandomService>.excludeRapidApiServices() = this.filter { excludeRapidApiService(it) }
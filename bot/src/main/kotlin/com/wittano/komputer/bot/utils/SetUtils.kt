package com.wittano.komputer.bot.utils

import com.wittano.komputer.bot.dagger.isDatabaseReady
import com.wittano.komputer.bot.joke.JokeApiService
import com.wittano.komputer.bot.joke.JokeCategory
import com.wittano.komputer.bot.joke.JokeRandomService
import com.wittano.komputer.bot.joke.JokeType
import com.wittano.komputer.bot.joke.api.rapidapi.RapidAPIService
import com.wittano.komputer.bot.joke.mongodb.JokeDatabaseService
import java.util.*

internal fun Set<JokeRandomService>.filterService(
    type: JokeType,
    category: JokeCategory,
    language: Locale = Locale.ENGLISH
): List<JokeRandomService> =
    this.filter {
        if (it is JokeApiService) {
            return@filter it.supports(type) && it.supports(category) && it.supports(language)
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
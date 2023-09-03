package com.wittano.komputer.core.utils

import com.wittano.komputer.core.config.dagger.isDatabaseReady
import com.wittano.komputer.core.joke.JokeApiService
import com.wittano.komputer.core.joke.JokeCategory
import com.wittano.komputer.core.joke.JokeRandomService
import com.wittano.komputer.core.joke.JokeType
import com.wittano.komputer.core.joke.api.rapidapi.RapidAPIService
import com.wittano.komputer.core.joke.mongodb.JokeDatabaseService

fun Set<JokeRandomService>.filterService(type: JokeType, category: JokeCategory): List<JokeRandomService> =
    this.filter {
        if (it is JokeApiService) {
            return@filter it.supports(type) && it.supports(category)
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

fun List<JokeRandomService>.excludeRapidApiServices() = this.filter { excludeRapidApiService(it) }
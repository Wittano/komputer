package com.wittano.komputer.utils

import com.wittano.komputer.config.dagger.isDatabaseReady
import com.wittano.komputer.joke.JokeApiService
import com.wittano.komputer.joke.JokeCategory
import com.wittano.komputer.joke.JokeRandomService
import com.wittano.komputer.joke.JokeType
import com.wittano.komputer.joke.api.rapidapi.RapidAPIService
import com.wittano.komputer.joke.mongodb.JokeDatabaseService

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
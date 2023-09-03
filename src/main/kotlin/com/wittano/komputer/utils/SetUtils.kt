package com.wittano.komputer.utils

import com.wittano.komputer.joke.JokeApiService
import com.wittano.komputer.joke.JokeCategory
import com.wittano.komputer.joke.JokeRandomService
import com.wittano.komputer.joke.JokeType
import com.wittano.komputer.joke.api.rapidapi.RapidAPIService

fun Set<JokeRandomService>.filterService(type: JokeType, category: JokeCategory): List<JokeRandomService> =
    this.filter {
        if (it is JokeApiService) {
            return@filter it.supports(type) && it.supports(category)
        }

        true
    }.filter {
        if (it is RapidAPIService) {
            return@filter it.isEnable()
        }

        true
    }
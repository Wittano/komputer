package com.wittano.komputer.bot.joke

import java.util.*

interface JokeApiService {

    fun supports(type: JokeType): Boolean = true
    fun supports(category: JokeCategory): Boolean

    fun supports(language: Locale) = language == Locale.ENGLISH

}
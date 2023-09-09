package com.wittano.komputer.bot.joke

import java.util.Locale

// TODO Add more external API integration
interface JokeApiService {

    fun supports(type: JokeType): Boolean = true
    fun supports(category: JokeCategory): Boolean

    fun supports(language: Locale) = language == Locale.ENGLISH

}
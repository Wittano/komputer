package com.wittano.komputer.bot.joke.api.rapidapi

import com.wittano.komputer.commons.config.config

interface RapidAPIService {
    fun isEnable(): Boolean = config.rapidApiKey?.isNotBlank() == true
}
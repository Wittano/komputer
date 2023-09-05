package com.wittano.komputer.core.joke.api.rapidapi

import com.wittano.komputer.core.config.config

interface RapidAPIService {
    fun isEnable(): Boolean = config.rapidApiKey?.isNotBlank() == true
}
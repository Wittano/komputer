package com.wittano.komputer.core.joke.api.rapidapi

import com.wittano.komputer.core.config.ConfigLoader

interface RapidAPIService {
    fun isEnable(): Boolean = ConfigLoader.load().rapidApiKey?.isNotBlank() == true
}
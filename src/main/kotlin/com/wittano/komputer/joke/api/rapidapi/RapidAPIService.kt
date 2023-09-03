package com.wittano.komputer.joke.api.rapidapi

import com.wittano.komputer.config.ConfigLoader

interface RapidAPIService {
    fun isEnable(): Boolean = ConfigLoader.load().rapidApiKey?.isNotBlank() == true
}
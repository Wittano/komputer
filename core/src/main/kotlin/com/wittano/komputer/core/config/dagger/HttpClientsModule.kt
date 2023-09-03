package com.wittano.komputer.core.config.dagger

import dagger.Provides
import okhttp3.OkHttpClient
import java.time.Duration
import javax.inject.Singleton

@Module
class HttpClientsModule {

    @Provides
    @Singleton
    fun createJokeDevHttpClient() = OkHttpClient()
        .newBuilder()
        .connectTimeout(Duration.ofSeconds(2))
        .readTimeout(Duration.ofSeconds(2))
        .build()
}

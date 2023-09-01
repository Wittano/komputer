package com.wittano.komputer.config.dagger

import dagger.Module
import dagger.Provides
import okhttp3.OkHttpClient
import java.time.Duration
import javax.inject.Named
import javax.inject.Singleton

@Module
class HttpClientsModule {

    @Provides
    @Named("jokeDevClient")
    @Singleton
    fun createJokeDevHttpClient() = OkHttpClient()
        .newBuilder()
        .connectTimeout(Duration.ofSeconds(2))
        .readTimeout(Duration.ofSeconds(1))
        .build()
}

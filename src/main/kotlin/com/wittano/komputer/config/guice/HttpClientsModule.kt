package com.wittano.komputer.config.guice

import com.google.inject.AbstractModule
import com.google.inject.Provides
import com.google.inject.name.Named
import okhttp3.OkHttpClient
import java.time.Duration

class HttpClientsModule : AbstractModule() {

    @Named("jokeDevClient")
    @Provides
    fun createJokeDevHttpClient() = OkHttpClient()
        .newBuilder()
        .connectTimeout(Duration.ofSeconds(2))
        .readTimeout(Duration.ofSeconds(1))
        .build()

}
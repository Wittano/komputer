package com.wittano.komputer.config.dagger

import com.fasterxml.jackson.databind.ObjectMapper
import com.mongodb.reactivestreams.client.MongoClient
import com.wittano.komputer.joke.JokeApiService
import com.wittano.komputer.joke.JokeRandomService
import com.wittano.komputer.joke.JokeService
import com.wittano.komputer.joke.jokedev.JokeDevClient
import com.wittano.komputer.joke.mongodb.JokeDatabaseService
import dagger.Module
import dagger.Provides
import dagger.multibindings.IntoSet
import okhttp3.OkHttpClient
import javax.inject.Inject
import javax.inject.Named
import javax.inject.Singleton

@Module
class JokeServiceModule {

    @Provides
    @Inject
    @Singleton
    fun createJokeDatabaseService(client: MongoClient): JokeService = createJokeDatabase(client)

    @Provides
    @Singleton
    @Inject
    fun createJokeDevService(
        objectMapper: ObjectMapper,
        @Named("jokeDevClient") client: OkHttpClient
    ): JokeApiService = createJokeDev(objectMapper, client)

    @Provides
    @Inject
    @IntoSet
    @Singleton
    fun createJokeRandomDatabaseService(client: MongoClient): JokeRandomService = createJokeDatabase(client)

    @Provides
    @Singleton
    @Inject
    @IntoSet
    fun createJokeRandomDevService(
        objectMapper: ObjectMapper,
        @Named("jokeDevClient") client: OkHttpClient
    ): JokeRandomService = createJokeDev(objectMapper, client)

    private fun createJokeDev(objectMapper: ObjectMapper, client: OkHttpClient) = JokeDevClient(client, objectMapper)

    private fun createJokeDatabase(client: MongoClient) = JokeDatabaseService(client)
}

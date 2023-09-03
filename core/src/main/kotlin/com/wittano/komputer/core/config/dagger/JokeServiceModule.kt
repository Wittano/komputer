package com.wittano.komputer.core.config.dagger

import com.fasterxml.jackson.databind.ObjectMapper
import com.mongodb.reactivestreams.client.MongoClient
import com.wittano.komputer.joke.JokeApiService
import com.wittano.komputer.joke.JokeRandomService
import com.wittano.komputer.joke.JokeService
import com.wittano.komputer.joke.api.jokedev.JokeDevClient
import com.wittano.komputer.joke.api.rapidapi.humorapi.HumorAPIService
import com.wittano.komputer.joke.mongodb.JokeDatabaseService
import dagger.Provides
import dagger.multibindings.IntoSet
import okhttp3.OkHttpClient
import javax.inject.Inject
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
        client: OkHttpClient
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
        client: OkHttpClient
    ): JokeRandomService = createJokeDev(objectMapper, client)

    @Provides
    @Inject
    @IntoSet
    @Singleton
    fun createJokeHumorAPIService(
        databaseService: JokeDatabaseService,
        client: OkHttpClient,
        objectMapper: ObjectMapper
    ): JokeApiService = createHumorAPIService(client, databaseService, objectMapper)

    @Provides
    @Singleton
    @Inject
    @IntoSet
    fun createJokeRandomHumorAPIService(
        databaseService: JokeDatabaseService,
        client: OkHttpClient,
        objectMapper: ObjectMapper
    ): JokeRandomService = createHumorAPIService(client, databaseService, objectMapper)

    private fun createJokeDev(objectMapper: ObjectMapper, client: OkHttpClient) = JokeDevClient(client, objectMapper)

    private fun createJokeDatabase(client: MongoClient) = JokeDatabaseService(client)

    private fun createHumorAPIService(
        client: OkHttpClient,
        database: JokeDatabaseService,
        objectMapper: ObjectMapper
    ) = HumorAPIService(client, database, objectMapper)
}

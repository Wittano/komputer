package com.wittano.komputer.bot.dagger

import com.fasterxml.jackson.databind.ObjectMapper
import com.wittano.komputer.bot.joke.JokeApiService
import com.wittano.komputer.bot.joke.JokeRandomService
import com.wittano.komputer.bot.joke.JokeService
import com.wittano.komputer.bot.joke.api.jokedev.JokeDevClient
import com.wittano.komputer.bot.joke.api.rapidapi.humorapi.HumorAPIService
import com.wittano.komputer.bot.joke.mongodb.JokeDatabaseService
import dagger.Module
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
    fun createJokeDatabaseService(): JokeService = createJokeDatabase()

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
    fun createJokeRandomDatabaseService(): JokeRandomService = createJokeDatabase()

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

    private fun createJokeDatabase() = JokeDatabaseService()

    private fun createHumorAPIService(
        client: OkHttpClient,
        database: JokeDatabaseService,
        objectMapper: ObjectMapper
    ) = HumorAPIService(client, database, objectMapper)
}

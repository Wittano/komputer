package com.wittano.komputer.config.dagger

import com.mongodb.reactivestreams.client.MongoClient
import com.wittano.komputer.joke.JokeService
import com.wittano.komputer.joke.mongodb.JokeDatabaseService
import dagger.Module
import dagger.Provides
import javax.inject.Inject
import javax.inject.Singleton

@Module
class JokeServiceModule {

    @Provides
    @Inject
    @Singleton
    fun createJokeDatabaseService(client: MongoClient): JokeService = JokeDatabaseService(client)

}

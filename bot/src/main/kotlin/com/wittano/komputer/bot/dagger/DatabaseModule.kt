package com.wittano.komputer.bot.dagger

import com.wittano.komputer.bot.config.ConfigDatabaseService
import com.wittano.komputer.bot.joke.mongodb.JokeDatabaseService
import dagger.Module
import dagger.Provides
import javax.inject.Singleton

@Module
class DatabaseModule {

    @Provides
    @Singleton
    fun createJokeDatabaseService() = JokeDatabaseService()

    @Provides
    @Singleton
    fun createConfigDatabaseService() = ConfigDatabaseService()

}
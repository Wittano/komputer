package com.wittano.komputer.bot.dagger

import com.wittano.komputer.bot.command.*
import com.wittano.komputer.bot.config.ConfigDatabaseService
import com.wittano.komputer.bot.joke.JokeRandomService
import com.wittano.komputer.bot.joke.api.jokedev.JokeDevClient
import com.wittano.komputer.bot.joke.mongodb.JokeDatabaseService
import dagger.Module
import dagger.Provides
import dagger.multibindings.IntoMap
import dagger.multibindings.StringKey
import javax.inject.Inject
import javax.inject.Singleton

@Module
class SlashCommandsModule {

    @Provides
    @StringKey("welcome")
    @IntoMap
    @Singleton
    fun createWelcomeCommand(): SlashCommand = WelcomeCommand()

    @Provides
    @StringKey("addjoke")
    @IntoMap
    @Inject
    @Singleton
    fun createAddJokeCommand(databaseManager: JokeDatabaseService): SlashCommand = AddJokeCommand(databaseManager)

    @Inject
    @IntoMap
    @Provides
    @StringKey("joke")
    @Singleton
    fun createJokeCommand(
        jokeDevClient: JokeDevClient,
        jokeRandomServices: Set<@JvmSuppressWildcards JokeRandomService>
    ): SlashCommand = JokeCommand(jokeDevClient, jokeRandomServices)

    @Inject
    @IntoMap
    @Provides
    @StringKey("removejoke")
    @Singleton
    fun createRemoveJokeCommand(service: JokeDatabaseService): SlashCommand = RemoveJokeCommand(service)

    @Inject
    @Provides
    @Singleton
    @IntoMap
    @StringKey("showconfig")
    fun createShowConfigCommand(configDatabaseService: ConfigDatabaseService): SlashCommand =
        ShowConfigCommand(configDatabaseService)

    @Inject
    @Provides
    @Singleton
    @IntoMap
    @StringKey("config")
    fun createUpdateConfigCommand(configDatabaseService: ConfigDatabaseService): SlashCommand =
        UpdateConfigCommand(configDatabaseService)
}
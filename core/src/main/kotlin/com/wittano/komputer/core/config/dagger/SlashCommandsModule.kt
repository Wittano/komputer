package com.wittano.komputer.core.config.dagger

import com.wittano.komputer.command.*
import com.wittano.komputer.core.command.WelcomeCommand
import com.wittano.komputer.joke.JokeRandomService
import com.wittano.komputer.joke.api.jokedev.JokeDevClient
import com.wittano.komputer.joke.mongodb.JokeDatabaseService
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

}
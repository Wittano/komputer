package com.wittano.komputer.config.dagger

import com.wittano.komputer.command.*
import com.wittano.komputer.joke.jokedev.JokeDevClient
import com.wittano.komputer.joke.mongodb.JokeDatabaseService
import dagger.Module
import dagger.Provides
import dagger.multibindings.IntoMap
import dagger.multibindings.StringKey
import javax.inject.Inject

@Module
class SlashCommandsModule {

    @Provides
    @StringKey("welcome")
    @IntoMap
    fun createWelcomeCommand(): SlashCommand = WelcomeCommand()

    @Provides
    @StringKey("addjoke")
    @IntoMap
    @Inject
    fun createAddJokeCommand(databaseManager: JokeDatabaseService): SlashCommand = AddJokeCommand(databaseManager)

    @Inject
    @IntoMap
    @Provides
    @StringKey("joke")
    fun createJokeCommand(jokeDevClient: JokeDevClient): SlashCommand = JokeCommand(jokeDevClient)

    @Inject
    @IntoMap
    @Provides
    @StringKey("removejoke")
    fun createRemoveJokeCommand(service: JokeDatabaseService): SlashCommand = RemoveJokeCommand(service)

}
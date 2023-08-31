package com.wittano.komputer.config.dagger

import com.wittano.komputer.command.AddJokeCommand
import com.wittano.komputer.command.JokeCommand
import com.wittano.komputer.command.SlashCommand
import com.wittano.komputer.command.WelcomeCommand
import com.wittano.komputer.joke.jokedev.JokeDevClient
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
    fun provideWelcomeCommand(): SlashCommand = WelcomeCommand()

    @Provides
    @StringKey("addjoke")
    @IntoMap
    fun provideAddJokeCommand(): SlashCommand = AddJokeCommand()

    @Inject
    @IntoMap
    @Provides
    @StringKey("joke")
    fun provideJokeCommand(jokeDevClient: JokeDevClient): SlashCommand = JokeCommand(jokeDevClient)

}
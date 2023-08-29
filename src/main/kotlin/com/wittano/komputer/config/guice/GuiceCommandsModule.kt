package com.wittano.komputer.config.guice

import com.google.inject.AbstractModule
import com.google.inject.Provides
import com.google.inject.name.Named
import com.wittano.komputer.command.AddJokeCommand
import com.wittano.komputer.command.JokeCommand
import com.wittano.komputer.command.SlashCommand
import com.wittano.komputer.command.WelcomeCommand

class GuiceCommandsModule : AbstractModule() {

    @Provides
    @Named("welcome")
    fun provideWelcomeCommand(): SlashCommand = WelcomeCommand()

    @Provides
    @Named("addjoke")
    fun provideAddJokeCommand(): SlashCommand = AddJokeCommand()

    @Provides
    @Named("add-joke")
    fun provideJokeCommand(): SlashCommand = JokeCommand()

}
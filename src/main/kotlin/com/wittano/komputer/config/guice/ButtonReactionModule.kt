package com.wittano.komputer.config.guice

import com.google.inject.AbstractModule
import com.google.inject.Inject
import com.google.inject.Provides
import com.google.inject.name.Named
import com.wittano.komputer.joke.jokedev.JokeDevClient
import com.wittano.komputer.message.interaction.ApologiesButtonReaction
import com.wittano.komputer.message.interaction.ButtonReaction
import com.wittano.komputer.message.interaction.NextJokeButtonReaction

class ButtonReactionModule : AbstractModule() {

    @Provides
    @Named("apologies")
    fun apologiesButton(): ButtonReaction = ApologiesButtonReaction()

    @Inject
    @Provides
    @Named("nextjoke")
    fun nextJokeButton(jokeDevClient: JokeDevClient): ButtonReaction = NextJokeButtonReaction(jokeDevClient)

}
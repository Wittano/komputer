package com.wittano.komputer.config.dagger

import com.wittano.komputer.joke.jokedev.JokeDevClient
import com.wittano.komputer.message.interaction.ApologiesButtonReaction
import com.wittano.komputer.message.interaction.ButtonReaction
import com.wittano.komputer.message.interaction.NextJokeButtonReaction
import dagger.Module
import dagger.Provides
import dagger.multibindings.IntoMap
import dagger.multibindings.StringKey
import javax.inject.Inject
import javax.inject.Singleton

@Module
class ButtonReactionModule {

    @Provides
    @IntoMap
    @StringKey("apologies")
    @Singleton
    fun apologiesButton(): ButtonReaction = ApologiesButtonReaction()

    @Inject
    @Provides
    @StringKey("nextjoke")
    @IntoMap
    @Singleton
    fun nextJokeButton(jokeDevClient: JokeDevClient): ButtonReaction = NextJokeButtonReaction(jokeDevClient)

}
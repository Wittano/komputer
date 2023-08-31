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

@Module
class ButtonReactionModule {

    @Provides
    @IntoMap
    @StringKey("apologies")
    fun apologiesButton(): ButtonReaction = ApologiesButtonReaction()

    @Inject
    @Provides
    @StringKey("nextjoke")
    @IntoMap
    fun nextJokeButton(jokeDevClient: JokeDevClient): ButtonReaction = NextJokeButtonReaction(jokeDevClient)

}
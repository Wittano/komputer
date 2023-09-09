package com.wittano.komputer.bot.dagger

import com.wittano.komputer.bot.joke.JokeRandomService
import com.wittano.komputer.bot.joke.api.jokedev.JokeDevClient
import com.wittano.komputer.bot.message.interaction.ApologiesButtonReaction
import com.wittano.komputer.bot.message.interaction.ButtonReaction
import com.wittano.komputer.bot.message.interaction.NextJokeButtonReaction
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
    fun nextJokeButton(
        jokeDevClient: JokeDevClient,
        jokeRandomServices: Set<@JvmSuppressWildcards JokeRandomService>
    ): ButtonReaction =
        NextJokeButtonReaction(jokeDevClient, jokeRandomServices)

}
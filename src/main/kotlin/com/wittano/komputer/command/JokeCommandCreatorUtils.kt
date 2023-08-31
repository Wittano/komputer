package com.wittano.komputer.command

import com.wittano.komputer.joke.JokeCategory
import com.wittano.komputer.joke.JokeType
import discord4j.discordjson.json.ApplicationCommandOptionChoiceData

internal val JOKE_CATEGORIES by lazy {
    JokeCategory.entries.map { it.toApplicationCommandOptionChoice() }
}

internal val JOKE_TYPES by lazy {
    JokeType.entries.map { it.toApplicationCommandOptionChoice() }
}

private fun JokeCategory.toApplicationCommandOptionChoice(): ApplicationCommandOptionChoiceData =
    ApplicationCommandOptionChoiceData.builder()
        .name(this.polishTranslate)
        .value(this.category)
        .build()

private fun JokeType.toApplicationCommandOptionChoice(): ApplicationCommandOptionChoiceData =
    ApplicationCommandOptionChoiceData.builder()
        .name(this.displayName)
        .value(this.jokeDevValue)
        .build()
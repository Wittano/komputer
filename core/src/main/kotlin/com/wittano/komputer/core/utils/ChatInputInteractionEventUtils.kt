package com.wittano.komputer.core.utils

import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import discord4j.core.`object`.command.ApplicationCommandInteractionOption
import kotlin.jvm.optionals.getOrNull

internal fun ChatInputInteractionEvent.getJokeCategory() =
    this.getOption("category")
        .flatMap(ApplicationCommandInteractionOption::getValue)
        .getOrNull()
        ?.asString()
        ?.let { category -> com.wittano.komputer.core.joke.JokeCategory.entries.find { it.category == category } }

internal fun ChatInputInteractionEvent.getJokeType() =
    this.getOption("type")
        .flatMap(ApplicationCommandInteractionOption::getValue)
        .getOrNull()
        ?.asString()
        ?.let { type -> com.wittano.komputer.core.joke.JokeType.entries.find { it.jokeDevValue == type } }
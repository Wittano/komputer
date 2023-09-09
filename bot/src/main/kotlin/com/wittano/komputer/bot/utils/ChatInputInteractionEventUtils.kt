package com.wittano.komputer.bot.utils

import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import discord4j.core.`object`.command.ApplicationCommandInteractionOption
import kotlin.jvm.optionals.getOrNull

internal fun ChatInputInteractionEvent.getJokeCategory() =
    this.getOption("category")
        .flatMap(ApplicationCommandInteractionOption::getValue)
        .getOrNull()
        ?.asString()
        ?.let { category -> com.wittano.komputer.bot.joke.JokeCategory.entries.find { it.category == category } }

internal fun ChatInputInteractionEvent.getJokeType() =
    this.getOption("type")
        .flatMap(ApplicationCommandInteractionOption::getValue)
        .getOrNull()
        ?.asString()
        ?.let { type -> com.wittano.komputer.bot.joke.JokeType.entries.find { it.jokeDevValue == type } }
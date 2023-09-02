package com.wittano.komputer.utils

import com.wittano.komputer.joke.JokeCategory
import com.wittano.komputer.joke.JokeType
import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import discord4j.core.`object`.command.ApplicationCommandInteractionOption
import kotlin.jvm.optionals.getOrNull

fun ChatInputInteractionEvent.getJokeCategory() =
    this.getOption("category")
        .flatMap(ApplicationCommandInteractionOption::getValue)
        .getOrNull()
        ?.asString()
        ?.let { category -> JokeCategory.entries.find { it.category == category } }

fun ChatInputInteractionEvent.getJokeType() =
    this.getOption("type")
        .flatMap(ApplicationCommandInteractionOption::getValue)
        .getOrNull()
        ?.asString()
        ?.let { type -> JokeType.entries.find { it.jokeDevValue == type } }
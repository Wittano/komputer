package com.wittano.komputer.bot.utils

import com.wittano.komputer.bot.joke.JokeCategory
import com.wittano.komputer.bot.joke.JokeType
import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import discord4j.core.`object`.command.ApplicationCommandInteractionOption
import java.util.*
import kotlin.jvm.optionals.getOrNull

internal fun ChatInputInteractionEvent.getJokeCategory() =
    this.getOption("category")
        .flatMap(ApplicationCommandInteractionOption::getValue)
        .getOrNull()
        ?.asString()
        ?.let { category -> JokeCategory.entries.find { it.category == category } }

internal fun ChatInputInteractionEvent.getJokeType() =
    this.getOption("type")
        .flatMap(ApplicationCommandInteractionOption::getValue)
        .getOrNull()
        ?.asString()
        ?.let { type -> JokeType.entries.find { it.type == type } }

// TODO Change default language based on global configuration per server
internal fun ChatInputInteractionEvent.getLanguage() =
    this.getOption("language")
        .flatMap(ApplicationCommandInteractionOption::getValue)
        .getOrNull()
        ?.asString()
        ?.let { Locale(it) }
        ?: Locale.ENGLISH
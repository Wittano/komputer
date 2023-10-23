package com.wittano.komputer.bot.utils.joke

import com.wittano.komputer.bot.joke.JokeCategory
import com.wittano.komputer.bot.joke.JokeType
import com.wittano.komputer.bot.utils.mongodb.getGlobalLanguage
import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import discord4j.core.event.domain.interaction.DeferrableInteractionEvent
import discord4j.core.`object`.command.ApplicationCommandInteractionOption
import discord4j.rest.util.Permission
import reactor.core.publisher.Mono
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

internal fun ChatInputInteractionEvent.getLanguage() =
    this.getOption("language")
        .flatMap(ApplicationCommandInteractionOption::getValue)
        .getOrNull()
        ?.asString()
        ?.let { Locale(it) }
        ?: getGlobalLanguage(this.getGuid())

internal fun ChatInputInteractionEvent.getLanguageOptional(): Locale? =
    this.getOption("language")
        .flatMap(ApplicationCommandInteractionOption::getValue)
        .getOrNull()
        ?.asString()
        ?.let { Locale(it) }

internal fun DeferrableInteractionEvent.getGuid(): String = this.interaction.guildId.get().asString()

internal fun DeferrableInteractionEvent.isAdministrator(): Mono<Boolean> =
    this.interaction.member.getOrNull()?.basePermissions?.map { permissionSet ->
        permissionSet.any { it == Permission.ADMINISTRATOR }
    } ?: Mono.just(false)
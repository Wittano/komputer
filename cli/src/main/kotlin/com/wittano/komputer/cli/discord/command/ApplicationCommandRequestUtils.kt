package com.wittano.komputer.cli.discord.command

import discord4j.discordjson.json.ApplicationCommandData
import discord4j.discordjson.json.ApplicationCommandOptionChoiceData
import discord4j.discordjson.json.ApplicationCommandOptionData
import discord4j.discordjson.json.ApplicationCommandRequest
import discord4j.discordjson.possible.Possible
import java.util.*
import kotlin.jvm.optionals.getOrNull

internal fun ApplicationCommandRequest.equalsCommand(command: ApplicationCommandData): Boolean {
    val requestOptions = this.options()
        .toOptional()
        .getOrNull()
        ?.apply {
            this.sortBy { it.name() }
        }

    val registeredCommandOptions = command.options()
        .toOptional()
        .getOrNull()
        ?.apply {
            this.sortBy { it.name() }
        }

    val isOptionsEquals = requestOptions == registeredCommandOptions
            || requestOptions?.zip(registeredCommandOptions ?: Collections.emptyList())
        ?.all { it.first.compareApplicationCommandOptions(it.second) } == true

    return this.name() == command.name()
            && this.description().toOptional().getOrNull() == command.description()
            && isOptionsEquals
}

private fun ApplicationCommandOptionData.compareApplicationCommandOptions(another: ApplicationCommandOptionData): Boolean {
    val anotherChoices = another.choices().getOrNull()
    val isChoicesEquals = this.choices().toOptional().getOrNull() == anotherChoices ||
            this.choices().getOrNull()
                ?.zip(anotherChoices ?: Collections.emptyList())?.all {
                    it.first.compareApplicationCommandOptionChoiceData(it.second)
                } == true

    val isChannelTypesEquals =
        this.channelTypes().getOrNull()?.toSet() == another.channelTypes().getOrNull()?.toSet()

    return this.name() == another.name() &&
            this.description() == another.description() &&
            this.required().toBoolean() == another.required().toBoolean() &&
            this.autocomplete().toBoolean() == another.autocomplete().toBoolean() &&
            this.type() == another.type() &&
            isChannelTypesEquals &&
            isChoicesEquals &&
            this.maxLength().getOrNull() == another.maxLength().getOrNull() &&
            this.minLength().getOrNull() == another.minLength().getOrNull() &&
            this.minValue().getOrNull() == another.minValue().getOrNull() &&
            this.maxValue().getOrNull() == another.maxValue().getOrNull()
}

private fun <T> Possible<out T>.getOrNull() = this.takeIf { !it.isAbsent }?.toOptional()?.get()

private fun Possible<Boolean>.toBoolean(): Boolean = this.toOptional().getOrNull() == true

private fun ApplicationCommandOptionChoiceData.compareApplicationCommandOptionChoiceData(another: ApplicationCommandOptionChoiceData): Boolean {
    return this.name() == another.name() && this.value() == another.value()
}

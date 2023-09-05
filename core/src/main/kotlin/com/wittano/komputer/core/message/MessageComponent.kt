package com.wittano.komputer.core.message

import com.wittano.komputer.core.message.resource.ButtonLabel
import com.wittano.komputer.core.message.resource.getButtonLabel
import discord4j.core.`object`.component.Button
import java.util.*

const val APOLOGIES_BUTTON_ID = "apologies"

const val NEXT_JOKE_BUTTON_ID = "next-joke"

internal fun createJokeReactionButtons(locale: Locale = Locale("pl")): List<Button> {
    val apologiesButton = Button.primary(APOLOGIES_BUTTON_ID, getButtonLabel(ButtonLabel.APOLOGIES, locale))
    val nextJoke = Button.secondary(NEXT_JOKE_BUTTON_ID, getButtonLabel(ButtonLabel.NEXT_JOKE, locale))

    return listOf(apologiesButton, nextJoke)
}


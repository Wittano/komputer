package com.wittano.komputer.core.message

import discord4j.core.`object`.component.Button

const val APOLOGIES_BUTTON_ID = "apologies"

const val NEXT_JOKE_BUTTON_ID = "next-joke"

// TODO Add internationalization button labels
internal fun createJokeReactionButtons(): List<Button> {
    val apologiesButton = createApologiesButton()
    val nextJoke = Button.secondary(NEXT_JOKE_BUTTON_ID, "Zabawne powiedz więcej")

    return listOf(apologiesButton, nextJoke)
}

internal fun createApologiesButton(): Button = Button.primary(APOLOGIES_BUTTON_ID, "Przeproś")
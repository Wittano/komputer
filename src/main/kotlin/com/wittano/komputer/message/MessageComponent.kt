package com.wittano.komputer.message

import discord4j.core.`object`.component.Button

const val APOLOGIES_BUTTON_ID = "apologies"

const val NEXT_JOKE_BUTTON_ID = "next-joke"

fun createJokeReactionButtons(): List<Button> {
    val apologiesButton = Button.primary(APOLOGIES_BUTTON_ID, "Przeproś")
    val nextJoke = Button.secondary(NEXT_JOKE_BUTTON_ID, "Zabawne powiedz więcej")

    return listOf(apologiesButton, nextJoke)
}
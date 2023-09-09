package com.wittano.komputer.core.message

import discord4j.core.spec.InteractionApplicationCommandCallbackSpec

internal fun createErrorMessage(): InteractionApplicationCommandCallbackSpec =
    InteractionApplicationCommandCallbackSpec.builder()
        // TODO Add english translate
        .content("BEEP BOOP. Coś poszło nie tak :(")
        .build()
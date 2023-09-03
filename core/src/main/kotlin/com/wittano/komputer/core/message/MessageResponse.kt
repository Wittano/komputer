package com.wittano.komputer.core.message

import discord4j.core.spec.InteractionApplicationCommandCallbackSpec

fun createErrorMessage() = InteractionApplicationCommandCallbackSpec.builder()
    .content("BEEP BOOP. Coś poszło nie tak :(")
    .build()
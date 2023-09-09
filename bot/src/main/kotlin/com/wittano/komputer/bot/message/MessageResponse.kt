package com.wittano.komputer.bot.message

import com.wittano.komputer.commons.extensions.POLISH_LOCALE
import com.wittano.komputer.commons.transtation.ErrorMessage
import com.wittano.komputer.commons.transtation.getErrorMessage
import discord4j.core.spec.InteractionApplicationCommandCallbackSpec
import java.util.*

internal fun createErrorMessage(locale: Locale = POLISH_LOCALE): InteractionApplicationCommandCallbackSpec =
    InteractionApplicationCommandCallbackSpec.builder()
        .content(getErrorMessage(ErrorMessage.GENERAL_ERROR, locale))
        .build()
package com.wittano.komputer.commons.transtation

import com.wittano.komputer.commons.extensions.POLISH_LOCALE
import java.util.*

fun getErrorMessage(errorMessage: ErrorMessage, locale: Locale = POLISH_LOCALE): String =
    getResourceBundle("i18n.error-message", locale).getString(errorMessage.code)

fun getButtonLabel(buttonLabel: ButtonLabel, locale: Locale = POLISH_LOCALE): String =
    getResourceBundle("i18n.reaction-button", locale).getString(buttonLabel.code)

fun getSuccessfulMessage(successfulMessage: SuccessfulMessage, locale: Locale = POLISH_LOCALE): String =
    getResourceBundle("i18n.successful-message", locale).getString(successfulMessage.code)

private fun getResourceBundle(path: String, locale: Locale) = ResourceBundle.getBundle(path, locale)
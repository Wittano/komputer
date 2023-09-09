package com.wittano.komputer.commons.transtation

import java.util.*

private val POLISH_LOCALE = Locale("pl")

fun getErrorMessage(code: ErrorMessage, locale: Locale = POLISH_LOCALE): String {
    val resourceBundle = ResourceBundle.getBundle("i18n.error-message", locale)

    return resourceBundle.getString(code.code)
}

fun getButtonLabel(name: ButtonLabel, locale: Locale = POLISH_LOCALE): String {
    val resource = ResourceBundle.getBundle("i18n.reaction-button", locale)

    return resource.getString(name.code)
}
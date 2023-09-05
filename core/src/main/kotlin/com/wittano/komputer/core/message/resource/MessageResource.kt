package com.wittano.komputer.core.message.resource

import java.util.*

internal fun getErrorMessage(code: ErrorMessage, locale: Locale = Locale.ENGLISH): String {
    val resourceBundle = ResourceBundle.getBundle("i18n.error-message", locale)

    return resourceBundle.getString(code.code)
}
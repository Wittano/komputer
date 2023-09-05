package com.wittano.komputer.core.message.resource

import java.util.*

fun getButtonLabel(name: ButtonLabel, locale: Locale): String {
    val resource = ResourceBundle.getBundle("i18n.reaction-button", locale)

    return resource.getString(name.code)
}
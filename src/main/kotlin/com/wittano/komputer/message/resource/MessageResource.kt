package com.wittano.komputer.message.resource

import java.util.*

class MessageResource private constructor() {

    companion object {
        fun get(code: ErrorMessage, locale: Locale = Locale.ENGLISH): String {
            val resourceBundle = ResourceBundle.getBundle("i18n.error-message", locale)

            return resourceBundle.getString(code.code)
        }
    }

}
package com.wittano.komputer.commons.extensions

import java.util.*

val POLISH_LOCALE = Locale("pl")

fun String.toLocale(): Locale {
    val localeString = this.split("-")

    return if (localeString.size == 2) {
        Locale(localeString[0], localeString[1])
    } else {
        Locale(this)
    }
}
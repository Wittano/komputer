package com.wittano.komputer.bot.joke

import com.wittano.komputer.commons.extensions.POLISH_LOCALE
import java.util.*

data class Joke(
    val answer: String,
    var category: JokeCategory,
    val type: JokeType,
    val question: String? = null,
    val language: Locale = POLISH_LOCALE
) {
    fun isYoMama(): Boolean {
        if (category == JokeCategory.YO_MAMA) {
            return true
        }

        val yoMana = ResourceBundle.getBundle("i18n.yo-mama-list", this.language)
            .getString("yo-mama.joke.prefix").split(",")

        return yoMana.any {
            this.answer.startsWith(it)
        }
    }
}

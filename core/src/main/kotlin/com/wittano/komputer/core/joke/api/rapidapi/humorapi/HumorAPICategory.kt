package com.wittano.komputer.core.joke.api.rapidapi.humorapi

import com.wittano.komputer.core.joke.JokeCategory

internal enum class HumorAPICategory(val category: JokeCategory, val tag: String) {
    YO_MAMA(JokeCategory.YO_MAMA, "yo_mama"),
    DARK(JokeCategory.DARK, "dark"),
    PROGRAMMING(JokeCategory.PROGRAMMING, "nerdy"),
    ONE_LINER(JokeCategory.ANY, "one_liner"),
}
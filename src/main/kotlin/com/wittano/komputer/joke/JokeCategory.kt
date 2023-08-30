package com.wittano.komputer.joke

enum class JokeCategory(val category: String, val polishTranslate: String) {
    PROGRAMMING("Programming", "Programowanie"),
    ANY("Any", "Dowolne"),
    MISC("Misc", "Misc"),
    DARK("Dark", "Czarny humor"),
    SPOOKY("Spooky", "Straszne")
}
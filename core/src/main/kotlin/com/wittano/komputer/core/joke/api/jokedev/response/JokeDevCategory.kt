package com.wittano.komputer.core.joke.api.jokedev.response

import com.fasterxml.jackson.core.JsonParser
import com.fasterxml.jackson.databind.DeserializationContext
import com.fasterxml.jackson.databind.JsonDeserializer
import com.fasterxml.jackson.databind.JsonNode
import com.fasterxml.jackson.databind.annotation.JsonDeserialize
import com.wittano.komputer.core.joke.JokeCategory

@JsonDeserialize(using = JokeDevCategoryDeserializer::class)
enum class JokeDevCategory(val value: String, val jokeCategory: JokeCategory) {
    PROGRAMMING("Programming", JokeCategory.PROGRAMMING),
    MISC("Misc", JokeCategory.MISC),
    DARK("Dark", JokeCategory.DARK),
    PUN("Pun", JokeCategory.ANY),
    SPOOKY("Spooky", JokeCategory.SPOOKY),
    CHRISTMAS("Christmas", JokeCategory.ANY)
}

class JokeDevCategoryDeserializer : JsonDeserializer<JokeDevCategory>() {
    override fun deserialize(jsonParser: JsonParser?, context: DeserializationContext?): JokeDevCategory? {
        val category = jsonParser?.codec?.readTree<JsonNode>(jsonParser)?.asText()

        return JokeDevCategory.entries.find { it.value == category }
    }

}
package com.wittano.komputer.message

import com.wittano.komputer.joke.Joke
import com.wittano.komputer.joke.JokeType
import discord4j.core.spec.EmbedCreateFields
import discord4j.core.spec.EmbedCreateSpec
import discord4j.rest.util.Color

fun createJokeMessage(joke: Joke): EmbedCreateSpec {
    val builder = EmbedCreateSpec.builder()
        .color(Color.of(0x02f5f5))
        .title("Joke")
        .author("komputer", null, null)


    if (joke.type == JokeType.TWO_PART) {
        val question = EmbedCreateFields.Field.of("Question", joke.question!!, false)
        val answer = EmbedCreateFields.Field.of("Answer", joke.answer, false)

        builder.addFields(question, answer)
    } else {
        builder.addField("Joke", joke.answer, false)
    }

    builder.addField("Category", joke.category.polishTranslate, false)

    return builder.build()
}
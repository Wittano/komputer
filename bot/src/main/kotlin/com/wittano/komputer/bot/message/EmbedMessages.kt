package com.wittano.komputer.bot.message

import com.wittano.komputer.bot.joke.Joke
import com.wittano.komputer.bot.joke.JokeCategory
import com.wittano.komputer.bot.joke.JokeType
import discord4j.core.spec.EmbedCreateFields
import discord4j.core.spec.EmbedCreateSpec
import discord4j.rest.util.Color

internal fun createJokeMessage(joke: Joke): EmbedCreateSpec {
    val builder = EmbedCreateSpec.builder()
        .color(Color.of(0x02f5f5))
        .title("Joke")
        .author("komputer", null, null) // TODO Add icon


    if (joke.type == JokeType.TWO_PART) {
        val question = EmbedCreateFields.Field.of("Question", joke.question!!, false)
        val answer = EmbedCreateFields.Field.of("Answer", joke.answer, false)

        builder.addFields(question, answer)
    } else {
        builder.addField("Joke", joke.answer, false)
    }

    if (joke.category == JokeCategory.YO_MAMA) {
        builder.image("https://media.tenor.com/sgS8GdoZGn8AAAAd/muscle-man-regular-show-muscle-man.gif")
    }

    builder.addField("Category", joke.category.polishTranslate, false)

    return builder.build()
}
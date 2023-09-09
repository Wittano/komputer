package com.wittano.komputer.bot.joke.mongodb

import com.wittano.komputer.bot.joke.JokeCategory
import com.wittano.komputer.bot.joke.JokeType
import org.bson.codecs.pojo.annotations.BsonId
import org.bson.codecs.pojo.annotations.BsonProperty
import org.bson.types.ObjectId

data class JokeModel(
    @BsonProperty("answer")
    val answer: String,
    @BsonProperty("type")
    val type: JokeType,
    @BsonProperty("category")
    val category: JokeCategory
) {
    @BsonId
    @BsonProperty("_id")
    lateinit var id: ObjectId

    @BsonProperty("question")
    var question: String? = null
}

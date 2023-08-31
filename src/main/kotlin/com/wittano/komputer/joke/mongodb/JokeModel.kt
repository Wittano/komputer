package com.wittano.komputer.joke.mongodb

import com.wittano.komputer.joke.JokeCategory
import com.wittano.komputer.joke.JokeType
import org.bson.codecs.pojo.annotations.BsonId
import org.bson.types.ObjectId

data class JokeModel(
    val answer: String,
    val question: String?,
    val type: JokeType,
    val category: JokeCategory
) {
    @BsonId
    lateinit var id: ObjectId
}

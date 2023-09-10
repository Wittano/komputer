package com.wittano.komputer.bot.config

import org.bson.codecs.pojo.annotations.BsonId
import org.bson.codecs.pojo.annotations.BsonProperty
import org.bson.types.ObjectId
import java.util.*

data class ServerConfigModel(
    @BsonProperty("guid")
    val guid: String,
    @BsonProperty("language")
    val language: String
) {
    @BsonId
    @BsonProperty("_id")
    lateinit var id: ObjectId

    fun toServerConfig() = ServerConfig(Locale(language))
}
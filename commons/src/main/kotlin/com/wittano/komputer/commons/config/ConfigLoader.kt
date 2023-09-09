package com.wittano.komputer.commons.config

import io.github.cdimascio.dotenv.Dotenv

data class Config(
    val token: String,
    val applicationId: Long,
    val guildId: Long,
    val mongoDbUri: String,
    val mongoDbName: String,
    val rapidApiKey: String?,
)

val config by lazy {
    val dotenv = Dotenv.configure().directory(System.getProperty("user.dir")).load()

    Config(
        token = dotenv.getOrElseThrow("DISCORD_BOT_TOKEN"),
        applicationId = dotenv.getOrElseThrow("APPLICATION_ID").toLong(),
        guildId = dotenv.getOrElseThrow("SERVER_GUID").toLong(),
        mongoDbUri = dotenv.getOrElseThrow("MONGODB_URI"),
        mongoDbName = dotenv.getOrElseThrow("MONGODB_DB_NAME"),
        rapidApiKey = dotenv.getOrNull("RAPID_API_KEY")
    )
}

private fun Dotenv.getOrNull(key: String): String? = this.get(key, System.getenv(key)).takeIf { it.isNotBlank() }

private fun Dotenv.getOrElseThrow(key: String): String = this.get(key, System.getenv(key)).takeIf { it.isNotBlank() }
    ?: throw IllegalStateException("Environment variable $key is missing")
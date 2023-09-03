package com.wittano.komputer.core.config

data class Config(
    val token: String,
    val applicationId: Long,
    val guildId: Long,
    val mongoDbUri: String,
    val mongoDbName: String,
    val rapidApiKey: String?,
)

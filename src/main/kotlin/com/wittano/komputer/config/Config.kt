package com.wittano.komputer.config

data class Config(
    val token: String,
    val applicationId: String,
    val serverGuid: String,
    val mongoDbUri: String,
)

package com.wittano.komputer.config

import io.github.cdimascio.dotenv.Dotenv

class ConfigLoader private constructor() {

    companion object {
        private var config: Config? = null

        fun load(): Config {
            if (config != null) {
                return config as Config
            }

            val dotenv = Dotenv.load()

            val loadedConfig = Config(
                token = dotenv.getOrElseThrow("DISCORD_BOT_TOKEN"),
                applicationId = dotenv.getOrElseThrow("APPLICATION_ID").toLong(),
                guildId = dotenv.getOrElseThrow("SERVER_GUID").toLong(),
                mongoDbUri = dotenv.getOrElseThrow("MONGODB_URI")
            )

            config = loadedConfig

            return loadedConfig
        }
    }

}

private fun Dotenv.getOrElseThrow(key: String): String = this[key].takeIf { it.isNotBlank() }
    ?: throw IllegalStateException("Environment variable $key is missing")
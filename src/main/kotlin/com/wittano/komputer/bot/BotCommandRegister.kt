package com.wittano.komputer.bot

import com.wittano.komputer.config.ConfigLoader
import discord4j.discordjson.json.ApplicationCommandData
import discord4j.discordjson.json.ApplicationCommandRequest
import discord4j.rest.RestClient
import reactor.core.publisher.Flux
import reactor.kotlin.core.publisher.toFlux

class BotCommandRegister private constructor() {

    companion object {
        fun registerCommands(
            client: RestClient,
            commands: List<ApplicationCommandRequest>,
            registeredCommands: Flux<ApplicationCommandData>
        ): Flux<ApplicationCommandData> {
            val config = ConfigLoader.load()

            return commands.toFlux()
                .filterWhen { request ->
                    registeredCommands.any { request.name().equals(it.name(), true) }
                }
                .flatMapSequential {
                    client.applicationService.createGuildApplicationCommand(config.applicationId, config.guildId, it)
                }
        }
    }

}
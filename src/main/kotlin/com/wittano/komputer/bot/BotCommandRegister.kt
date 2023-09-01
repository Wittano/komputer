package com.wittano.komputer.bot

import com.wittano.komputer.config.ConfigLoader
import discord4j.discordjson.json.ApplicationCommandData
import discord4j.discordjson.json.ApplicationCommandRequest
import discord4j.rest.RestClient
import reactor.core.publisher.Flux
import reactor.core.publisher.Mono

internal class BotCommandRegister private constructor() {

    companion object {
        fun registerCommands(
            client: RestClient,
            commands: List<ApplicationCommandRequest>,
            registeredCommands: Mono<MutableList<ApplicationCommandData>>
        ): Flux<ApplicationCommandData> {
            val config = ConfigLoader.load()

            return Flux.fromIterable(commands)
                .filterWhen { request ->
                    registeredCommands.map {
                        it.none { data -> data.name().equals(request.name(), true) }
                    }.switchIfEmpty(Mono.just(true))
                }
                .flatMapSequential {
                    client.applicationService.createGuildApplicationCommand(config.applicationId, config.guildId, it)
                }
        }
    }

}
package com.wittano.komputer.bot

import com.wittano.komputer.config.ConfigLoader
import discord4j.discordjson.json.ApplicationCommandData
import discord4j.discordjson.json.ApplicationCommandRequest
import discord4j.rest.RestClient
import reactor.core.publisher.Flux
import reactor.core.publisher.Mono

internal class BotCommandCleaner private constructor() {
    companion object {
        fun deleteUnusedGuildCommands(
            client: RestClient,
            commands: MutableList<ApplicationCommandRequest>,
            registeredCommands: Mono<MutableList<ApplicationCommandData>>
        ): Flux<in Unit> {
            val config = ConfigLoader.load()

            return registeredCommands.flatMapMany {
                Flux.fromIterable(it)
            }.filter {
                !commands.map(ApplicationCommandRequest::name).contains(it.name())
            }.flatMap {
                client.applicationService.deleteGuildApplicationCommand(
                    config.applicationId,
                    config.guildId,
                    it.id().asLong()
                )
            }
        }
    }


}
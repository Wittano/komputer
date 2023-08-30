package com.wittano.komputer.bot

import com.wittano.komputer.config.ConfigLoader
import discord4j.discordjson.json.ApplicationCommandData
import discord4j.discordjson.json.ApplicationCommandRequest
import discord4j.rest.RestClient
import reactor.core.publisher.Flux

class BotCommandCleaner private constructor() {
    companion object {
        fun deleteUnusedGuildCommands(
            client: RestClient,
            commands: MutableList<ApplicationCommandRequest>,
            registeredCommands: Flux<ApplicationCommandData>
        ): Flux<in Unit> {
            val config = ConfigLoader.load()

            return registeredCommands.filter {
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
package com.wittano.komputer.cli.command

import com.wittano.komputer.bot.discordClient
import com.wittano.komputer.cli.discord.command.RegisteredCommandsUtils
import com.wittano.komputer.commons.config.config
import org.slf4j.LoggerFactory
import picocli.CommandLine.Command
import reactor.core.publisher.Flux

@Command(
    name = "init",
    description = ["Register slash commands on specified server or global"]
)
class BotInitCommands : Runnable {
    private val log = LoggerFactory.getLogger(this::class.qualifiedName)

    // TODO Add parameter to sing in commands as global

    override fun run() {
        val commands = RegisteredCommandsUtils.getCommandsFromJsonFiles()

        Flux.fromIterable(commands)
            .flatMap {
                discordClient.restClient.applicationService.createGuildApplicationCommand(
                    config.applicationId,
                    config.guildId,
                    it
                ).doOnSuccess { commandData ->
                    log.info("Successfully added new command: '${commandData.name()}'")
                }
            }.doOnComplete {
                log.info("Komputer's commands was created successfully")

            }.collectList().block()

        discordClient.logout().block()
    }
}
package com.wittano.komputer.cli.command

import com.wittano.komputer.bot.discordClient
import com.wittano.komputer.cli.discord.DiscordException
import com.wittano.komputer.cli.discord.command.RegisteredCommandsUtils
import com.wittano.komputer.cli.discord.command.equalsCommand
import com.wittano.komputer.commons.config.config
import discord4j.discordjson.json.ApplicationCommandData
import discord4j.discordjson.json.ApplicationCommandRequest
import org.slf4j.LoggerFactory
import picocli.CommandLine.Command
import picocli.CommandLine.Option
import reactor.core.publisher.Flux
import reactor.core.publisher.Mono

@Command(
    name = "update",
    description = ["Update slash command for specified server"]
)
class BotCommandsUpdater : Runnable {

    private val log = LoggerFactory.getLogger(this::class.qualifiedName)

    @Option(
        names = ["--enable-global"],
        description = ["Flag if program should register command as global"]
    )
    private var isSingInGlobal: Boolean = false

    override fun run() {
        val commands = RegisteredCommandsUtils.getCommandsFromJsonFiles()

        updateCommands(commands).toIterable().forEach { command ->
            val isCommandEqual = commands.any {
                it.equalsCommand(command)
            }

            check(isCommandEqual) {
                "Command '${command.name()}' didn't update correctly'"
            }
        }

        log.info("Komputer's commands updated successfully")

        discordClient.logout().block()
    }

    private fun updateCommands(commands: List<ApplicationCommandRequest>): Flux<ApplicationCommandData> {
        val updateApplicationCommands = if (isSingInGlobal) {
            discordClient.restClient.applicationService.bulkOverwriteGlobalApplicationCommand(
                config.applicationId,
                commands
            )
        } else {
            discordClient.restClient.applicationService.bulkOverwriteGuildApplicationCommand(
                config.applicationId,
                config.guildId,
                commands
            )
        }

        return updateApplicationCommands.switchIfEmpty(Mono.error(DiscordException("Failed update command")))
    }
}
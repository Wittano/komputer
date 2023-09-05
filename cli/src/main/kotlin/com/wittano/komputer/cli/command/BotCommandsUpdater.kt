package com.wittano.komputer.cli.command

import com.wittano.komputer.cli.discord.DiscordException
import com.wittano.komputer.cli.discord.command.RegisteredCommandsUtils
import com.wittano.komputer.cli.discord.command.equalsCommand
import com.wittano.komputer.core.bot.discordClient
import com.wittano.komputer.core.config.config
import discord4j.discordjson.json.ApplicationCommandData
import discord4j.discordjson.json.ApplicationCommandRequest
import org.slf4j.LoggerFactory
import picocli.CommandLine.Command
import picocli.CommandLine.Parameters
import reactor.core.publisher.Flux
import reactor.core.publisher.Mono

@Command(
    name = "update",
    description = ["Update slash command for specified server"]
)
class BotCommandsUpdater : Runnable {

    private val log = LoggerFactory.getLogger(this::class.qualifiedName)

    @Parameters(description = ["List of command's names, which will be updated"])
    lateinit var commandName: Array<String>

    override fun run() {
        if (commandName.isEmpty()) {
            log.warn("Command's names list is empty")
            return
        }

        val commands = RegisteredCommandsUtils.getCommandsFromJsonFiles()
            .filter {
                commandName.contains(it.name())
            }

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
        return discordClient.restClient.applicationService.bulkOverwriteGuildApplicationCommand(
            config.applicationId,
            config.guildId,
            commands
        ).switchIfEmpty { Mono.error<ApplicationCommandData>(DiscordException("Failed update command")) }
    }
}
package com.wittano.komputer.cli.command

import com.wittano.komputer.core.bot.discordClient
import com.wittano.komputer.core.config.config
import org.slf4j.LoggerFactory
import picocli.CommandLine.Command
import picocli.CommandLine.Parameters

@Command(
    name = "remove",
    description = ["Remove command from specified server or global scope"]
)
class BotCommandRemover : Runnable {
    private val log = LoggerFactory.getLogger(this::class.qualifiedName)

    @Parameters(description = ["List of command's names, which will be removed"])
    private lateinit var commandNames: Array<String>

    override fun run() {
        if (commandNames.isEmpty()) {
            log.warn("List of commands to remove is empty")
            return
        }

        val commands = discordClient.restClient
            .applicationService.getGuildApplicationCommands(
                config.applicationId,
                config.guildId
            ).filter {
                commandNames.contains(it.name())
            }

        commands.flatMap {
            discordClient.restClient.applicationService.deleteGuildApplicationCommand(
                config.applicationId,
                config.guildId,
                it.id().asLong()
            )
        }.doOnComplete {
            log.info("Komputer's commands '${commandNames.joinToString()}' removed successfully")
        }.collectList().block()

        discordClient.logout().block()
    }

}
package com.wittano.komputer.cli

import com.wittano.komputer.cli.command.BotCommandRemover
import com.wittano.komputer.cli.command.BotCommandsUpdater
import com.wittano.komputer.cli.command.BotInitCommands
import com.wittano.komputer.cli.command.BotRunner
import org.slf4j.LoggerFactory
import picocli.CommandLine
import kotlin.system.exitProcess

fun main(args: Array<String>) {
    val logger = LoggerFactory.getLogger(Thread.currentThread().name)

    try {
        CommandLine(BotRunner())
            .addSubcommand(BotCommandsUpdater())
            .addSubcommand(BotCommandRemover())
            .addSubcommand(BotInitCommands())
            .execute(*args)
    } catch (ex: Exception) {
        logger.error("Unhandled exception: ${ex.message}")
        exitProcess(-1)
    }
}
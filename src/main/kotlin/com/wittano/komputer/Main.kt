package com.wittano.komputer

import com.google.inject.Guice
import com.wittano.komputer.bot.KomputerBot
import com.wittano.komputer.config.guice.GuiceCommandsModule
import org.slf4j.LoggerFactory
import picocli.CommandLine
import kotlin.system.exitProcess

fun main(args: Array<String>) {
    val logger = LoggerFactory.getLogger("MAIN")

    try {
        val injector = Guice.createInjector(GuiceCommandsModule())

        CommandLine(KomputerBot(injector)).execute(*args)
    } catch (ex: Exception) {
        logger.error("Unhandled exception: ${ex.message}")
        exitProcess(-1)
    }
}
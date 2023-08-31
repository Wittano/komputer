package com.wittano.komputer

import com.wittano.komputer.bot.KomputerBot
import org.slf4j.LoggerFactory
import picocli.CommandLine
import kotlin.system.exitProcess

fun main(args: Array<String>) {
    val logger = LoggerFactory.getLogger(Thread.currentThread().name)

    try {
        CommandLine(KomputerBot()).execute(*args)
    } catch (ex: Exception) {
        logger.error("Unhandled exception: ${ex.message}")
        exitProcess(-1)
    }
}
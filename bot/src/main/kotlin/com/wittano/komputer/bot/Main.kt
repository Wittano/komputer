package com.wittano.komputer.bot

import com.wittano.komputer.bot.bot.KomputerBot
import org.slf4j.LoggerFactory
import kotlin.system.exitProcess

fun main() {
    val logger = LoggerFactory.getLogger(Thread.currentThread().name)

    try {
        KomputerBot().start()
    } catch (ex: Exception) {
        logger.error("Unhandled exception: ${ex.message}", ex)
        exitProcess(-1)
    }
}
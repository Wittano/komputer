package com.wittano.komputer.core

import com.wittano.komputer.core.bot.KomputerBot
import org.slf4j.LoggerFactory
import kotlin.system.exitProcess

// TODO Make internal scope for function and variables
fun main() {
    val logger = LoggerFactory.getLogger(Thread.currentThread().name)

    try {
        KomputerBot().start()
    } catch (ex: Exception) {
        logger.error("Unhandled exception: ${ex.message}", ex)
        exitProcess(-1)
    }
}
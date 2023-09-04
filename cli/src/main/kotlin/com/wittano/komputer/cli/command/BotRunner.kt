package com.wittano.komputer.cli.command

import com.wittano.komputer.core.bot.KomputerBot
import picocli.CommandLine.Command

@Command(
    name = "komputer",
    description = ["Discord bot behave as like \"komputer\". One of character in Star Track parody series created by Dem3000"],
)
class BotRunner : Runnable {
    override fun run() {
        KomputerBot().start()
    }
}
package com.wittano.komputer.command

import discord4j.discordjson.json.ApplicationCommandRequest

interface Command {

    fun execute()

    fun createCommand(): ApplicationCommandRequest

}
package com.wittano.komputer.bot

import discord4j.common.JacksonResources
import discord4j.discordjson.json.ApplicationCommandRequest
import discord4j.rest.RestClient
import discord4j.rest.interaction.GlobalCommandRegistrar
import java.nio.file.FileSystems
import java.nio.file.Files
import java.nio.file.Path
import java.nio.file.Paths

class BotCommandRegister(private val client: RestClient) {

    fun registerCommands() {
        val jacksonResources = JacksonResources.create()
        val commandDirectoryPath = getCommandDirectory()
        val commands = mutableListOf<ApplicationCommandRequest>()

        commandDirectoryPath.toFile().listFiles()?.forEach {
            val commandConfig = Files.readAllBytes(it.toPath())
            val command = jacksonResources.objectMapper.readValue(commandConfig, ApplicationCommandRequest::class.java)

            commands.add(command)
        }

        GlobalCommandRegistrar.create(client, commands).registerCommands().blockFirst()
    }

    private fun getCommandDirectory(): Path {
        val uri = this::class.java.classLoader?.getResource("commands")?.toURI()

        return if ("jar" == uri?.scheme) {
            val fileSystem = FileSystems.newFileSystem(uri, mutableMapOf<String, Any>())

            fileSystem.getPath("src/main/resources/commands")
        } else {
            uri?.let { Paths.get(it) } ?: throw IllegalStateException("Commands directory wasn't found")
        }
    }

}
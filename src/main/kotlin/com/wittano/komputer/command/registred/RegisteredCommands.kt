package com.wittano.komputer.command.registred

import discord4j.common.JacksonResources
import discord4j.discordjson.json.ApplicationCommandRequest
import java.nio.file.FileSystems
import java.nio.file.Files
import java.nio.file.Path
import java.nio.file.Paths

class RegisteredCommandsUtils private constructor() {
    companion object {
        fun getCommandsFromJsonFiles(): MutableList<ApplicationCommandRequest> {
            val jacksonResources = JacksonResources.create()
            val commands = mutableListOf<ApplicationCommandRequest>()

            getCommandDirectory().toFile()
                .listFiles()
                ?.forEach {
                    val commandConfig = Files.readAllBytes(it.toPath())
                    val command =
                        jacksonResources.objectMapper.readValue(commandConfig, ApplicationCommandRequest::class.java)

                    commands.add(command)
                }

            return commands
        }

        private fun getCommandDirectory(): Path {
            val uri = RegisteredCommandsUtils::class.java.classLoader?.getResource("commands")?.toURI()

            return if ("jar" == uri?.scheme) {
                val fileSystem = FileSystems.newFileSystem(uri, mutableMapOf<String, Any>())

                fileSystem.getPath("src/main/resources/commands")
            } else {
                uri?.let { Paths.get(it) } ?: throw IllegalStateException("Commands directory wasn't found")
            }
        }
    }

}

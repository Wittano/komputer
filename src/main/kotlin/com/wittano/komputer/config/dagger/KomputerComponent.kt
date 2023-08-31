package com.wittano.komputer.config.dagger

import com.wittano.komputer.command.SlashCommand
import com.wittano.komputer.message.interaction.ButtonReaction
import dagger.Component

@Component(
    modules = [
        ButtonReactionModule::class,
        HttpClientsModule::class,
        SlashCommandsModule::class,
        UtilitiesModule::class,
        MongoDbModule::class
    ]
)
interface KomputerComponent {
    fun getSlashCommand(): Map<String, SlashCommand>
    fun getButtonReaction(): Map<String, ButtonReaction>
}
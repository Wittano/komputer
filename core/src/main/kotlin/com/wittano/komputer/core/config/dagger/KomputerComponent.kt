package com.wittano.komputer.core.config.dagger

import com.wittano.komputer.core.command.SlashCommand
import com.wittano.komputer.core.message.interaction.ButtonReaction
import dagger.Component
import javax.inject.Singleton

@Component(
    modules = [
        ButtonReactionModule::class,
        HttpClientsModule::class,
        SlashCommandsModule::class,
        UtilitiesModule::class,
        JokeServiceModule::class,
        MongoDbModule::class,
    ]
)
@Singleton
interface KomputerComponent {
    fun getSlashCommand(): Map<String, SlashCommand>
    fun getButtonReaction(): Map<String, ButtonReaction>
}
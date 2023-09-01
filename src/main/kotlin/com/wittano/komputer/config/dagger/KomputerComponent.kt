package com.wittano.komputer.config.dagger

import com.wittano.komputer.command.SlashCommand
import com.wittano.komputer.message.interaction.ButtonReaction
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
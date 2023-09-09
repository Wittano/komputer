package com.wittano.komputer.bot.dagger

import com.wittano.komputer.bot.command.SlashCommand
import com.wittano.komputer.bot.message.interaction.ButtonReaction
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
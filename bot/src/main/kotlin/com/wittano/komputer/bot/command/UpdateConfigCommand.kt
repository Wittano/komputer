package com.wittano.komputer.bot.command

import com.wittano.komputer.bot.config.ConfigDatabaseService
import com.wittano.komputer.bot.config.ServerConfig
import com.wittano.komputer.bot.message.createErrorMessage
import com.wittano.komputer.bot.message.createSuccessfulMessage
import com.wittano.komputer.bot.utils.joke.getGuid
import com.wittano.komputer.bot.utils.mongodb.getGlobalLanguage
import com.wittano.komputer.commons.transtation.ErrorMessage
import com.wittano.komputer.commons.transtation.SuccessfulMessage
import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import discord4j.core.`object`.command.ApplicationCommandInteractionOption
import reactor.core.publisher.Mono
import java.util.*
import javax.inject.Inject
import kotlin.jvm.optionals.getOrNull

class UpdateConfigCommand @Inject constructor(
    private val configDatabaseService: ConfigDatabaseService
) : SlashCommand {
    override fun execute(event: ChatInputInteractionEvent): Mono<Void> {
        val guid = event.getGuid()
        val language = event.getOption("language")
            .flatMap(ApplicationCommandInteractionOption::getValue)
            .getOrNull()
            ?.asString()
            ?.let { Locale(it) }

        val updateMessage = event.reply(
            createSuccessfulMessage(
                SuccessfulMessage.CONFIG_UPDATED,
                language ?: getGlobalLanguage(guid)
            )
        )

        if (language == null) {
            return updateMessage
        }

        val config = ServerConfig(language)

        return configDatabaseService.update(event.getGuid(), config)
            .transform { updateMessage }
            .switchIfEmpty(event.reply(createErrorMessage(ErrorMessage.CONFIG_UPDATE_FAILED, language)))
    }
}
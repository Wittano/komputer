package com.wittano.komputer.bot.command

import com.wittano.komputer.bot.config.ConfigDatabaseService
import com.wittano.komputer.bot.config.ServerConfig
import com.wittano.komputer.bot.utils.joke.getGuid
import com.wittano.komputer.commons.transtation.ConfigFieldTranslation
import com.wittano.komputer.commons.transtation.getConfigFieldName
import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import discord4j.core.spec.EmbedCreateFields
import discord4j.core.spec.EmbedCreateSpec
import discord4j.core.spec.InteractionApplicationCommandCallbackSpec
import reactor.core.publisher.Mono
import javax.inject.Inject

class ShowConfigCommand @Inject constructor(
    private val configDatabaseService: ConfigDatabaseService
) : SlashCommand {
    override fun execute(event: ChatInputInteractionEvent): Mono<Void> {
        val guid = event.getGuid()

        return configDatabaseService[guid]
            .flatMap {
                event.reply(createShowMessage(it))
            }
    }

    private fun createShowMessage(config: ServerConfig): InteractionApplicationCommandCallbackSpec {
        val builder = InteractionApplicationCommandCallbackSpec.builder()

        builder.addEmbed(
            EmbedCreateSpec.builder()
                .title(getConfigFieldName(ConfigFieldTranslation.TITLE, config.language))
                .addField(
                    createConfigField(
                        getConfigFieldName(ConfigFieldTranslation.LANGUAGE, config.language),
                        config.language.language
                    )
                )
                .build()
        )

        return builder.build()
    }

    private fun createConfigField(name: String, value: String): EmbedCreateFields.Field =
        EmbedCreateFields.Field.of(name, value, false)
}
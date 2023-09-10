package com.wittano.komputer.bot.command

import com.wittano.komputer.bot.command.exception.AccessDeniedException
import com.wittano.komputer.bot.config.ConfigDatabaseService
import com.wittano.komputer.bot.config.ServerConfig
import com.wittano.komputer.bot.utils.joke.getGuid
import com.wittano.komputer.bot.utils.joke.isAdministrator
import com.wittano.komputer.commons.transtation.ConfigFieldTranslation
import com.wittano.komputer.commons.transtation.getConfigFieldName
import discord4j.core.event.domain.interaction.ChatInputInteractionEvent
import discord4j.core.spec.EmbedCreateFields
import discord4j.core.spec.EmbedCreateSpec
import discord4j.core.spec.InteractionApplicationCommandCallbackSpec
import reactor.core.publisher.Mono
import javax.inject.Inject
import kotlin.jvm.optionals.getOrNull

class ShowConfigCommand @Inject constructor(
    private val configDatabaseService: ConfigDatabaseService
) : SlashCommand {
    override fun execute(event: ChatInputInteractionEvent): Mono<Void> {
        val isAdministrator = event.isAdministrator()
        val guid = event.getGuid()

        return isAdministrator.flatMap {
            if (it) {
                configDatabaseService[guid]
                    .flatMap { config ->
                        event.reply(createShowMessage(event, config))
                    }
            } else {
                val userId = event.interaction.user.id.asString()
                val guildId = event.getGuid()

                Mono.error(AccessDeniedException(userId, guildId))
            }
        }

    }

    private fun createShowMessage(
        event: ChatInputInteractionEvent,
        config: ServerConfig
    ): InteractionApplicationCommandCallbackSpec {
        val builder = InteractionApplicationCommandCallbackSpec.builder()
        val role = event.interaction.member.getOrNull()
            ?.roles
            ?.filter { it.id.asString() == config.roleId }
            ?.map { it.name }
            ?.blockFirst()

        builder.addEmbed(
            EmbedCreateSpec.builder()
                .title(getConfigFieldName(ConfigFieldTranslation.TITLE, config.language))
                .addFields(
                    createConfigField(
                        getConfigFieldName(ConfigFieldTranslation.LANGUAGE, config.language),
                        config.language.language
                    ),
                    createConfigField(
                        getConfigFieldName(ConfigFieldTranslation.ROLE, config.language),
                        role.orEmpty(),
                    )
                )
                .build()
        )

        return builder.build().withEphemeral(true)
    }

    private fun createConfigField(name: String, value: String): EmbedCreateFields.Field =
        EmbedCreateFields.Field.of(name, value, false)
}
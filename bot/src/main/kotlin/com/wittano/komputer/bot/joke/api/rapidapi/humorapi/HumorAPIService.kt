package com.wittano.komputer.bot.joke.api.rapidapi.humorapi

import com.fasterxml.jackson.databind.ObjectMapper
import com.wittano.komputer.bot.joke.*
import com.wittano.komputer.bot.joke.api.JokeApiException
import com.wittano.komputer.bot.joke.api.rapidapi.RapidAPIService
import com.wittano.komputer.bot.joke.api.rapidapi.RapidApiException
import com.wittano.komputer.bot.joke.mongodb.JokeDatabaseService
import com.wittano.komputer.commons.config.config
import com.wittano.komputer.commons.transtation.ErrorMessage
import okhttp3.HttpUrl.Companion.toHttpUrl
import okhttp3.OkHttpClient
import okhttp3.Request
import org.slf4j.LoggerFactory
import reactor.core.publisher.Mono
import reactor.core.scheduler.Schedulers
import reactor.kotlin.core.publisher.switchIfEmpty
import java.util.*
import java.util.concurrent.atomic.AtomicBoolean
import javax.inject.Inject

class HumorAPIService @Inject constructor(
    private val client: OkHttpClient,
    private val database: JokeDatabaseService,
    private val objectMapper: ObjectMapper
) : JokeApiService, JokeRandomService, RapidAPIService {

    private val log = LoggerFactory.getLogger(this::class.qualifiedName)
    private val isLimitExceeded = AtomicBoolean(false)

    override fun isEnable(): Boolean = super.isEnable() && !isLimitExceeded.get()

    override fun supports(type: JokeType): Boolean = type == JokeType.SINGLE

    override fun supports(category: JokeCategory): Boolean =
        category != JokeCategory.MISC && category != JokeCategory.SPOOKY

    override fun getRandom(category: JokeCategory?, type: JokeType, language: Locale?): Mono<Joke> {
        val humorCategory = category?.toHumorAPICategory() ?: HumorAPICategory.ONE_LINER
        val url = "https://humor-jokes-and-memes.p.rapidapi.com/jokes/random".toHttpUrl()
            .newBuilder()
            .addQueryParameter("exclude-tags", "nsfw")
            .addQueryParameter("include-tags", humorCategory.tag)
            .build()

        val request = Request.Builder()
            .url(url)
            .header("X-RapidAPI-Key", config.rapidApiKey!!)
            .header("X-RapidAPI-Host", "humor-jokes-and-memes.p.rapidapi.com")
            .get()
            .build()

        return Mono.just(client.newCall(request).execute())
            .flatMap {
                if (it.code == 429) {
                    isLimitExceeded.set(true)
                    return@flatMap Mono.error(RapidApiException("Limit of request for HumorAPI was exceeded"))
                }

                Mono.just(it)
            }
            .publishOn(Schedulers.boundedElastic())
            .mapNotNull<HumorAPIJokeResponse> {
                if (it.body == null) {
                    return@mapNotNull null
                }

                val responseStream = it.body?.byteStream() ?: return@mapNotNull null
                val response = objectMapper.readValue(responseStream, HumorAPIJokeResponse::class.java)
                responseStream.close()

                response
            }.map {
                Joke(
                    answer = it.joke,
                    category = humorCategory.category,
                    type = type
                )
            }.switchIfEmpty(
                Mono.error(
                    JokeApiException(
                        "Joke with type '${type}' and category '${humorCategory}' wasn't found",
                        ErrorMessage.JOKE_NOT_FOUND
                    )
                )
            ).publishOn(Schedulers.boundedElastic())
            .doOnSuccess {
                database.add(it)
                    .switchIfEmpty {
                        log.warn("Joke from HumorAPI is exist in database")

                        Mono.empty()
                    }.subscribe()
            }
    }

}

private fun JokeCategory.toHumorAPICategory(): HumorAPICategory =
    HumorAPICategory.entries.find { this == it.category }
        ?: HumorAPICategory.ONE_LINER
package com.wittano.komputer.joke.api.rapidapi.humorapi

import com.fasterxml.jackson.databind.ObjectMapper
import com.wittano.komputer.config.ConfigLoader
import com.wittano.komputer.joke.*
import com.wittano.komputer.joke.api.JokeApiException
import com.wittano.komputer.joke.api.rapidapi.RapidAPIService
import com.wittano.komputer.joke.mongodb.JokeDatabaseService
import com.wittano.komputer.message.resource.ErrorMessage
import okhttp3.HttpUrl.Companion.toHttpUrl
import okhttp3.OkHttpClient
import okhttp3.Request
import reactor.core.publisher.Mono
import reactor.core.scheduler.Schedulers
import javax.inject.Inject

class HumorAPIService @Inject constructor(
    private val client: OkHttpClient,
    private val database: JokeDatabaseService,
    private val objectMapper: ObjectMapper
) : JokeApiService, JokeRandomService, RapidAPIService {

    override fun supports(category: JokeCategory): Boolean =
        category != JokeCategory.MISC && category != JokeCategory.SPOOKY

    override fun getRandom(category: JokeCategory?, type: JokeType): Mono<Joke> {
        val humorCategory = category?.toHumorAPICategory() ?: HumorAPICategory.ONE_LINER
        val url = "https://humor-jokes-and-memes.p.rapidapi.com/jokes/random".toHttpUrl()
            .newBuilder()
            .addQueryParameter("exclude-tags", "nsfw")
            .addQueryParameter("include-tags", humorCategory.tag)
            .build()

        val request = Request.Builder()
            .url(url)
            .header("X-RapidAPI-Key", ConfigLoader.load().rapidApiKey!!)
            .header("X-RapidAPI-Host", "humor-jokes-and-memes.p.rapidapi.com")
            .get()
            .build()

        return Mono.just(client.newCall(request).execute())
            .flatMap {
                if (!it.isSuccessful || it.body == null) {
                    return@flatMap Mono.error(
                        JokeApiException(
                            "HumorAPI request failed. Response status ${it.code}, Body: ${it.body?.string()}",
                            ErrorMessage.JOKE_NOT_FOUND
                        )
                    )
                }

                Mono.just(it)
            }.mapNotNull<HumorAPIJokeResponse> {
                if (it.body == null) {
                    return@mapNotNull null
                }

                objectMapper.readValue(it.body?.bytes(), HumorAPIJokeResponse::class.java)
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
            .doOnSuccess { database.add(it).subscribe() }
    }

}

private fun JokeCategory.toHumorAPICategory(): HumorAPICategory =
    HumorAPICategory.entries.find { this == it.category }
        ?: HumorAPICategory.ONE_LINER
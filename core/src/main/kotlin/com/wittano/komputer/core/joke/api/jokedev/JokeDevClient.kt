package com.wittano.komputer.core.joke.api.jokedev

import com.fasterxml.jackson.databind.ObjectMapper
import com.wittano.komputer.core.joke.*
import com.wittano.komputer.core.joke.api.jokedev.response.JokeDevErrorResponse
import com.wittano.komputer.core.joke.api.jokedev.response.JokeDevSingleResponse
import com.wittano.komputer.core.joke.api.jokedev.response.JokeDevTwoPartResponse
import com.wittano.komputer.core.message.resource.ErrorMessage
import okhttp3.HttpUrl.Companion.toHttpUrl
import okhttp3.OkHttpClient
import okhttp3.Request
import org.slf4j.LoggerFactory
import reactor.core.publisher.Mono
import java.io.IOException
import javax.inject.Inject

class JokeDevClient @Inject constructor(
    private val client: OkHttpClient,
    private val objectMapper: ObjectMapper
) : JokeApiService, JokeRandomService {

    private val log = LoggerFactory.getLogger(this::class.qualifiedName)

    override fun getRandom(category: JokeCategory?, type: JokeType): Mono<Joke> {
        val apiCategory = category?.category ?: JokeCategory.ANY.category
        val requestUrl = "https://v2.jokeapi.dev/joke/${apiCategory}".toHttpUrl().newBuilder()
            .addQueryParameter("type", type.jokeDevValue)
            .build()

        val request = Request.Builder().url(requestUrl).build()

        val rawResponse = Mono.just(client.newCall(request).execute())
            .flatMap {
                if (!it.isSuccessful || it.body == null) {
                    return@flatMap Mono.error(
                        JokeDevApiException(
                            "JokeDev API request failed. Response status ${it.code}, Body: ${it.body?.string()}",
                            ErrorMessage.JOKE_NOT_FOUND
                        )
                    )
                }

                Mono.just(it)
            }

        val responseType: Class<out JokeExtractor> = if (type == JokeType.SINGLE) {
            JokeDevSingleResponse::class.java
        } else {
            JokeDevTwoPartResponse::class.java
        }

        val responseBytes = rawResponse.map { it.body!!.bytes() }

        return responseBytes.flatMap {
            try {
                Mono.just(objectMapper.readValue(it, responseType))
            } catch (ex: IOException) {
                val response = objectMapper.readValue(it, JokeDevErrorResponse::class.java)

                Mono.error(JokeDevApiException("Failed to get joke", ErrorMessage.JOKE_NOT_FOUND, response))
            }
        }.map {
            it.toJoke().apply {
                if (this.isYoMama()) {
                    this.category = JokeCategory.YO_MAMA
                }
            }
        }.doOnError {
            val errorResponse = (it as JokeDevApiException).response

            log.error("Failed get random joke from URL ${requestUrl}. Error message: ${errorResponse?.message}", it)
        }
    }

    override fun supports(category: JokeCategory): Boolean = category != JokeCategory.YO_MAMA

}
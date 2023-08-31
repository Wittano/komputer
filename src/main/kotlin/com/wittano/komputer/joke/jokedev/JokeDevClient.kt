package com.wittano.komputer.joke.jokedev

import com.fasterxml.jackson.databind.ObjectMapper
import com.wittano.komputer.joke.Joke
import com.wittano.komputer.joke.JokeCategory
import com.wittano.komputer.joke.JokeExtractor
import com.wittano.komputer.joke.JokeType
import com.wittano.komputer.joke.jokedev.response.JokeDevErrorResponse
import com.wittano.komputer.joke.jokedev.response.JokeDevSingleResponse
import com.wittano.komputer.joke.jokedev.response.JokeDevTwoPartResponse
import okhttp3.OkHttpClient
import okhttp3.Request
import org.slf4j.LoggerFactory
import java.io.IOException
import javax.inject.Inject
import javax.inject.Named

class JokeDevClient @Inject constructor(
    @Named("jokeDevClient")
    private val client: OkHttpClient,
    private val objectMapper: ObjectMapper
) {

    private val log = LoggerFactory.getLogger(this::class.qualifiedName)

    @Throws(JokeDevApiException::class)
    fun getRandomJoke(category: JokeCategory, type: JokeType): Joke {
        val requestUrl = "https://v2.jokeapi.dev/joke/${category.category}?type=${type.jokeDevValue}"
        val request = Request.Builder().url(requestUrl).build()

        val rawResponse = client.newCall(request).execute()
        if (!rawResponse.isSuccessful || rawResponse.body == null) {
            throw JokeDevApiException("JokeDev API request failed. Response status ${rawResponse.code}, Body: ${rawResponse.body?.string()}")
        }

        val responseType: Class<out JokeExtractor> = if (type == JokeType.SINGLE) {
            JokeDevSingleResponse::class.java
        } else {
            JokeDevTwoPartResponse::class.java
        }

        val responseBytes = rawResponse.body!!.bytes()
        val jokeResponse = try {
            objectMapper.readValue(responseBytes, responseType)
        } catch (ex: IOException) {
            objectMapper.readValue(responseBytes, JokeDevErrorResponse::class.java)
        }

        return when (jokeResponse) {
            is JokeDevErrorResponse -> {
                log.warn("Failed get random joke from URL ${requestUrl}. Error message: ${jokeResponse.message}")
                throw JokeDevApiException(jokeResponse.message)
            }

            else -> (jokeResponse as JokeExtractor).toJoke()
        }
    }

}
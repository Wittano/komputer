package com.wittano.komputer.joke.jokedev

import com.fasterxml.jackson.databind.ObjectMapper
import com.google.inject.Inject
import com.google.inject.name.Named
import com.wittano.komputer.joke.Joke
import com.wittano.komputer.joke.JokeCategory
import com.wittano.komputer.joke.JokeExtractor
import com.wittano.komputer.joke.JokeType
import com.wittano.komputer.joke.jokedev.response.JokeDevSingleResponse
import com.wittano.komputer.joke.jokedev.response.JokeDevTwoPartResponse
import okhttp3.OkHttpClient
import okhttp3.Request

class JokeDevClient @Inject constructor(
    @Named("jokeDevClient")
    private val client: OkHttpClient,
    private val objectMapper: ObjectMapper
) {

    fun getRandomJoke(category: JokeCategory, type: JokeType): Joke {
        val request = Request.Builder()
            .url("https://v2.jokeapi.dev/joke/${category.category}?type=${type.value}")
            .build()

        val rawResponse = client.newCall(request).execute()
        if (!rawResponse.isSuccessful || rawResponse.body == null) {
            throw JokeDevApiException("JokeDev API request failed. Response status ${rawResponse.code}, Body: ${rawResponse.body?.string()}")
        }

        val responseType: Class<out JokeExtractor> = if (type == JokeType.SINGLE) {
            JokeDevSingleResponse::class.java
        } else {
            JokeDevTwoPartResponse::class.java
        }

        val jokeResponse = objectMapper.readValue(rawResponse.body!!.byteStream(), responseType)

        return jokeResponse.toJoke()
    }

}
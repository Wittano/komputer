package com.wittano.komputer.core.joke.api.rapidapi

import com.wittano.komputer.core.joke.api.JokeApiException
import com.wittano.komputer.message.resource.ErrorMessage

class RapidApiException(msg: String) : JokeApiException(msg, ErrorMessage.JOKE_NOT_FOUND)
package com.wittano.komputer.bot.joke.api.rapidapi

import com.wittano.komputer.bot.joke.api.JokeApiException
import com.wittano.komputer.commons.transtation.ErrorMessage

class RapidApiException(msg: String) : JokeApiException(msg, ErrorMessage.JOKE_NOT_FOUND)
package com.wittano.komputer.bot.joke.mongodb

import com.mongodb.BasicDBObject
import com.mongodb.reactivestreams.client.MongoCollection
import com.wittano.komputer.bot.command.exception.CommandException
import com.wittano.komputer.bot.joke.*
import com.wittano.komputer.bot.utils.mongodb.MongoCollectionManager
import com.wittano.komputer.bot.utils.mongodb.database
import com.wittano.komputer.bot.utils.mongodb.getCollection
import com.wittano.komputer.bot.utils.mongodb.getModelCollection
import com.wittano.komputer.commons.transtation.ErrorMessage
import org.bson.Document
import org.bson.conversions.Bson
import org.bson.types.ObjectId
import org.slf4j.LoggerFactory
import reactor.core.publisher.Mono
import reactor.kotlin.core.publisher.toMono
import java.util.*

private const val JOKES_COLLECTION_NAME = "jokes"

class JokeDatabaseService : JokeService, JokeRandomService {
    private val log = LoggerFactory.getLogger(this::class.qualifiedName)

    init {
        database.flatMapMany {
            MongoCollectionManager.createCollectionIfDontExist(it, "jokes")
        }.subscribe()
    }

    override fun add(joke: Joke): Mono<String> {
        val jokeCollection = getModelCollection<JokeModel>(JOKES_COLLECTION_NAME)
        val isJokeAdded = jokeCollection.flatMap {
            val contentFilter = BasicDBObject().apply {
                this["answer"] = joke.answer
            }

            Mono.from(it.find(contentFilter))
                .map { true }
                .switchIfEmpty(Mono.just(false))
        }

        return jokeCollection
            .filterWhen { !isJokeAdded }
            .flatMap {
                it.insertOne(joke.toModel()).toMono()
            }.filter {
                it.wasAcknowledged()
            }.doOnError {
                log.error("Failed to add new joke into database. Cause: ${it.message}", it)
            }.mapNotNull {
                it.insertedId?.asObjectId()?.value?.toString()
            }
    }

    override fun remove(id: String): Mono<Void> {
        val jokeCollection = getCollection(JOKES_COLLECTION_NAME)

        return jokeCollection.flatMap {
            Mono.from(it.deleteOne(id.toBson()))
        }.doOnError {
            log.error("Failed to remove joke with id $id into database. Cause: ${it.message}", it)
        }.toMonoVoid()
    }

    override fun get(id: String): Mono<Joke> {
        val jokeCollection = getModelCollection<JokeModel>("")

        return jokeCollection.flatMap {
            try {
                Mono.from(it.find(id.toBson()))
            } catch (_: Exception) {
                Mono.error(InvalidJokeIdException("Joke ID is invalid", ErrorMessage.JOKE_ID_INVALID))
            }
        }.map {
            it.toJoke()
        }
    }

    override fun getRandom(category: JokeCategory?, type: JokeType?, language: Locale?): Mono<Joke> {
        val jokeCollection = getCollection(JOKES_COLLECTION_NAME)

        return jokeCollection.flatMap {
            findRandomJoke(it, category, type, language?.language)
        }.map {
            it.toJoke()
        }
    }

    private fun findRandomJoke(
        collection: MongoCollection<Document>,
        category: JokeCategory?,
        type: JokeType?,
        language: String?
    ): Mono<JokeModel> {
        val sampleObject = BasicDBObject().apply {
            this["\$sample"] = BasicDBObject().also { doc ->
                doc["size"] = 10
            }
        }

        val matcherObject = BasicDBObject().apply {
            this["\$match"] = BasicDBObject().apply {
                category?.takeIf { it != JokeCategory.ANY }
                    ?.also { c -> this["category"] = c }

                language?.also {
                    this["language"] = it
                }

                type?.also {
                    this["type"] = it.toString()
                }
            }
        }

        return Mono.from(collection.aggregate(mutableListOf(sampleObject, matcherObject), JokeModel::class.java))
            .switchIfEmpty(
                Mono.error(
                    CommandException(
                        "Joke with type '${type}', category '${category}' and language '$language' wasn't found",
                        ErrorMessage.JOKE_NOT_FOUND
                    )
                )
            )
    }
}

private operator fun Mono<Boolean>.not(): Mono<Boolean> = this.map { !it }

private fun Joke.toModel(): JokeModel = JokeModel(answer, type, category, language.language, question)

private fun String.toBson(): Bson = BasicDBObject().apply {
    this["_id"] = ObjectId(this@toBson)
}
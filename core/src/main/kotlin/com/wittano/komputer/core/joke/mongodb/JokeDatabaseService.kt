package com.wittano.komputer.core.joke.mongodb

import com.mongodb.BasicDBObject
import com.mongodb.reactivestreams.client.MongoClient
import com.mongodb.reactivestreams.client.MongoCollection
import com.mongodb.reactivestreams.client.MongoDatabase
import com.wittano.komputer.core.config.config
import com.wittano.komputer.core.joke.*
import com.wittano.komputer.core.message.resource.ErrorMessage
import org.bson.Document
import org.bson.conversions.Bson
import org.bson.types.ObjectId
import org.slf4j.LoggerFactory
import reactor.core.publisher.Flux
import reactor.core.publisher.Mono
import reactor.kotlin.core.publisher.toMono
import javax.inject.Inject

private const val JOKES_DATABASE_NAME = "jokes"

class JokeDatabaseService @Inject constructor(
    private val client: MongoClient
) : JokeService, JokeRandomService {
    private val log = LoggerFactory.getLogger(this::class.qualifiedName)
    private val database by lazy {
        val dbName = config.mongoDbName

        Mono.just(this.client)
            .map { it.getDatabase(dbName) }
            .doOnError { log.error("Failed to find '${dbName}' database", it) }
    }

    init {
        database.flatMap { createJokesCollectionIfDontExist(it) }
            .subscribe()
    }

    override fun add(joke: Joke): Mono<String> {
        val jokeCollection = getJokeModelCollection()
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
        val jokeCollection = getJokeCollection()

        return jokeCollection.flatMap {
            Mono.from(it.deleteOne(id.toBson()))
        }.doOnError {
            log.error("Failed to remove joke with id $id into database. Cause: ${it.message}", it)
        }.toMonoVoid()
    }

    override fun get(id: String): Mono<Joke> {
        val jokeCollection = getJokeModelCollection()

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

    override fun getRandom(category: JokeCategory?, type: JokeType): Mono<Joke> {
        val jokeCollection = getJokeCollection()

        return jokeCollection.flatMap {
            findRandomJoke(it, category, type)
        }.map {
            it.toJoke()
        }
    }

    private fun findRandomJoke(
        collection: MongoCollection<Document>,
        category: JokeCategory?,
        type: JokeType
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

                this["type"] = type.toString()
            }
        }

        return Mono.from(collection.aggregate(mutableListOf(sampleObject, matcherObject), JokeModel::class.java))
            .switchIfEmpty(
                Mono.error(
                    JokeException(
                        "Joke with type '${type}' and category '${category}' wasn't found",
                        ErrorMessage.JOKE_NOT_FOUND
                    )
                )
            )
    }

    private fun getJokeCollection() =
        database.map {
            it.getCollection(JOKES_DATABASE_NAME)
        }.doOnError {
            log.error("Failed get '$JOKES_DATABASE_NAME' collection", it)
        }

    private fun getJokeModelCollection() =
        database.map {
            it.getCollection(JOKES_DATABASE_NAME, JokeModel::class.java)
        }.doOnError {
            log.error("Failed get '$JOKES_DATABASE_NAME' collection", it)
        }

    // TODO Change function to pass collections names.
    // In the future, Komputer's database will have many collections e.g. config, jokes etc.
    private fun createJokesCollectionIfDontExist(database: MongoDatabase): Mono<Void> {
        val collectionsNames = Flux.from(database.listCollectionNames())

        return collectionsNames.collectList()
            .filter {
                !it.contains(JOKES_DATABASE_NAME)
            }.flatMap {
                Mono.from(database.createCollection(JOKES_DATABASE_NAME))
            }.doOnError {
                log.error("Failed to create '$JOKES_DATABASE_NAME' collection", it)
            }
    }

}

private operator fun Mono<Boolean>.not(): Mono<Boolean> = this.map { !it }

private fun Joke.toModel(): JokeModel =
    JokeModel(answer, type, category).apply { this.question = this@toModel.question }

private fun String.toBson(): Bson = BasicDBObject().apply {
    this["_id"] = ObjectId(this@toBson)
}

private fun JokeModel.toJoke(): Joke = Joke(answer, category, type, question)
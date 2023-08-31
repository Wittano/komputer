package com.wittano.komputer.joke.mongodb

import com.mongodb.BasicDBObject
import com.mongodb.reactivestreams.client.MongoClient
import com.mongodb.reactivestreams.client.MongoCollection
import com.mongodb.reactivestreams.client.MongoDatabase
import com.wittano.komputer.config.ConfigLoader
import com.wittano.komputer.joke.Joke
import com.wittano.komputer.joke.JokeCategory
import com.wittano.komputer.joke.JokeService
import com.wittano.komputer.joke.JokeType
import org.bson.Document
import org.bson.conversions.Bson
import org.bson.types.ObjectId
import org.slf4j.LoggerFactory
import reactor.core.publisher.Flux
import reactor.core.publisher.Mono
import reactor.kotlin.core.publisher.toMono
import javax.inject.Inject

private const val JOKES_DATABASE_NAME = "jokes"

class JokeDatabaseManager @Inject constructor(
    private val client: MongoClient
) : JokeService {
    private val log = LoggerFactory.getLogger(this::class.qualifiedName)
    private val database by lazy {
        val dbName = ConfigLoader.load().mongoDbName

        Mono.just(this.client)
            .map {
                it.getDatabase(dbName)
            }
            .doOnError {
                log.error("Failed to find '${dbName}' database", it)
            }
    }

    init {
        database.flatMap {
            createJokesCollectionIfDontExist(it)
        }.subscribe()
    }

    override fun add(joke: Joke): Mono<String> {
        val jokeCollection = getJokeCollection()

        return jokeCollection.flatMap {
            it.insertOne(joke.toDocument()).toMono()
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

        return jokeCollection
            .flatMap {
                Mono.from(it.deleteOne(id.toBson()))
            }.filter {
                it.wasAcknowledged()
            }.doOnError {
                log.error("Failed to remove joke with id $id into database. Cause: ${it.message}", it)
            }.toMonoVoid()
    }

    override fun get(id: String): Mono<Joke> {
        val jokeCollection = getJokeModelCollection()

        return jokeCollection.flatMap {
            Mono.from(it.find(id.toBson()))
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
        val sampleObject = BasicDBObject()
        sampleObject["\$sample"] = BasicDBObject().also { doc ->
            doc["size"] = 1
        }

        val matcherObject = BasicDBObject()
        matcherObject["\$match"] = BasicDBObject().apply {
            category?.also { c -> this["category"] = c }

            this["type"] = type.toString()
        }

        return Mono.from(collection.aggregate(mutableListOf(sampleObject, matcherObject), JokeModel::class.java))
    }

    private fun getJokeCollection() =
        database.map {
            it.getCollection(JOKES_DATABASE_NAME)
        }.doOnError {
            log.error("Failed get '${JOKES_DATABASE_NAME}' collection", it)
        }

    private fun getJokeModelCollection() =
        database.map {
            it.getCollection(JOKES_DATABASE_NAME, JokeModel::class.java)
        }.doOnError {
            log.error("Failed get '${JOKES_DATABASE_NAME}' collection", it)
        }

    private fun createJokesCollectionIfDontExist(database: MongoDatabase): Mono<Void> {
        val collectionsNames = Flux.from(database.listCollectionNames())
        val createCollection = Mono.from(database.createCollection(JOKES_DATABASE_NAME))

        return collectionsNames
            .filter {
                it == JOKES_DATABASE_NAME
            }
            .singleOrEmpty()
            .toMonoVoid()
            .switchIfEmpty(createCollection)
            .doOnError {
                log.error("Failed to create '${JOKES_DATABASE_NAME}' collection", it)
            }
    }

}

private fun String.toBson(): Bson {
    val bson = BasicDBObject()
    bson["_id"] = ObjectId(this)

    return bson
}

private fun JokeModel.toJoke(): Joke = Joke(answer, category, type, question)

private fun Joke.toDocument(): Document = Document(
    mapOf(
        Pair("category", this.category),
        Pair("type", this.type),
        Pair("question", this.question),
        Pair("answer", this.answer)
    )
)
package com.wittano.komputer.joke.mongodb

import com.mongodb.reactivestreams.client.MongoClient
import com.wittano.komputer.config.ConfigLoader
import com.wittano.komputer.joke.Joke
import org.bson.Document
import org.slf4j.LoggerFactory
import reactor.core.publisher.Flux
import reactor.core.publisher.Mono
import reactor.kotlin.core.publisher.toFlux
import reactor.kotlin.core.publisher.toMono
import javax.inject.Inject

private const val JOKES_DATABASE_NAME = "jokes"

class JokeDatabaseManager @Inject constructor(
    private val client: MongoClient
) {
    private val log = LoggerFactory.getLogger(this::class.qualifiedName)
    private val database by lazy {
        val dbName = ConfigLoader.load().mongoDbName

        Mono.just(client)
            .map {
                it.getDatabase(dbName)
            }
            .doOnError {
                log.error("Failed to find '${dbName}' database", it)
            }
    }

    fun addJoke(joke: Joke): Mono<String> {
        val jokeCollection = database.map {
            it.getCollection(JOKES_DATABASE_NAME)
        }.doOnError {
            log.error("Failed get '${JOKES_DATABASE_NAME}' collection", it)
        }

        return createJokesCollectionIfDontExist()
            .then(jokeCollection)
            .flatMap {
                it.insertOne(joke.toDocument()).toMono()
            }.filter {
                it.wasAcknowledged()
            }.doOnError {
                log.error("Failed to add new joke into database. Cause: ${it.message}", it)
            }.mapNotNull {
                it.insertedId?.asObjectId()?.value?.toString()
            }
    }

    private fun createJokesCollectionIfDontExist(): Mono<Void> {
        val collectionsNames = database.toFlux().flatMap {
            Flux.from(it.listCollectionNames())
        }

        val createCollection = database.flatMap {
            Mono.from(it.createCollection(JOKES_DATABASE_NAME))
        }.then(Mono.just(String()))

        return collectionsNames
            .filter {
                it == JOKES_DATABASE_NAME
            }
            .singleOrEmpty()
            .switchIfEmpty(createCollection)
            .mapNotNull<Void> { null }
            .doOnError {
                log.error("Failed to create '${JOKES_DATABASE_NAME}' collection", it)
            }
    }

}

private fun Joke.toDocument(): Document = Document(
    mapOf(
        Pair("category", this.category),
        Pair("type", this.type),
        Pair("question", this.question),
        Pair("answer", this.answer)
    )
)
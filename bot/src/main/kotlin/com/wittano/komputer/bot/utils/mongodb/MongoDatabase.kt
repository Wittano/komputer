package com.wittano.komputer.bot.utils.mongodb

import com.mongodb.ConnectionString
import com.mongodb.MongoClientSettings
import com.mongodb.ServerApi
import com.mongodb.ServerApiVersion
import com.mongodb.reactivestreams.client.MongoClient
import com.mongodb.reactivestreams.client.MongoClients
import com.wittano.komputer.commons.config.config
import org.bson.Document
import org.bson.codecs.configuration.CodecRegistries
import org.bson.codecs.pojo.PojoCodecProvider
import org.slf4j.LoggerFactory
import reactor.core.publisher.Mono
import java.time.Duration
import java.util.concurrent.atomic.AtomicBoolean

internal val database by lazy {
    val mongoDbName = config.mongoDbName

    Mono.just(client)
        .map { it.getDatabase(mongoDbName) }
        .doOnError { log.error("Failed to find '$mongoDbName' database", it) }
}

internal fun getCollection(name: String) =
    database.map {
        it.getCollection(name)
    }.doOnError {
        log.error("Failed get '$name' collection", it)
    }

internal inline fun <reified T> getModelCollection(name: String) =
    database.map {
        it.getCollection(name, T::class.java)
    }.doOnError {
        log.error("Failed get '$name' collection", it)
    }

val isDatabaseReady = AtomicBoolean(false)
private val log = LoggerFactory.getLogger("DATABASE_INITIALIZE")

private val client: MongoClient by lazy {
    val serverApi = ServerApi.builder()
        .version(ServerApiVersion.V1)
        .build()

    val pojoProvider = CodecRegistries.fromProviders(
        PojoCodecProvider.builder()
            .automatic(true)
            .build()
    )

    val codecRegister = CodecRegistries.fromRegistries(MongoClients.getDefaultCodecRegistry(), pojoProvider)

    val settings = MongoClientSettings.builder()
        .applyConnectionString(ConnectionString(config.mongoDbUri))
        .serverApi(serverApi)
        .codecRegistry(codecRegister)
        .build()

    val client = MongoClients.create(settings)

    checkDatabaseConnection(client)

    client
}

private fun checkDatabaseConnection(client: MongoClient) {
    Mono.from(client.getDatabase(config.mongoDbName).runCommand(Document("ping", 1)))
        .timeout(Duration.ofSeconds(2))
        .doOnError {
            log.error("Failed to connect with MongoDB database", it)
        }.doOnSuccess {
            isDatabaseReady.set(true)
        }.subscribe()
}
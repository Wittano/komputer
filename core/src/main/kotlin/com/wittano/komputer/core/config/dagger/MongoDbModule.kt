package com.wittano.komputer.core.config.dagger

import com.mongodb.ConnectionString
import com.mongodb.MongoClientSettings
import com.mongodb.ServerApi
import com.mongodb.ServerApiVersion
import com.mongodb.reactivestreams.client.MongoClient
import com.mongodb.reactivestreams.client.MongoClients
import com.wittano.komputer.core.config.Config
import com.wittano.komputer.core.config.ConfigLoader
import dagger.Module
import dagger.Provides
import org.bson.Document
import org.bson.codecs.configuration.CodecRegistries
import org.bson.codecs.pojo.PojoCodecProvider
import org.slf4j.LoggerFactory
import reactor.core.publisher.Mono
import java.time.Duration
import java.util.concurrent.atomic.AtomicBoolean
import javax.inject.Singleton

var isDatabaseReady = AtomicBoolean(false)

@Module
class MongoDbModule {

    private val log = LoggerFactory.getLogger(this::class.qualifiedName)

    @Provides
    @Singleton
    fun client(): MongoClient {
        val config = ConfigLoader.load()

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

        checkDatabaseConnection(client, config)

        return client
    }

    private fun checkDatabaseConnection(client: MongoClient, config: Config) {
        Mono.from(client.getDatabase(config.mongoDbName).runCommand(Document("ping", 1)))
            .timeout(Duration.ofSeconds(2))
            .doOnError {
                log.error("Failed to connect with MongoDB database", it)
            }.doOnSuccess {
                isDatabaseReady.set(true)
            }.subscribe()
    }

}
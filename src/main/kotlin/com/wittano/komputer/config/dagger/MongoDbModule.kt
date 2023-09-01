package com.wittano.komputer.config.dagger

import com.mongodb.ConnectionString
import com.mongodb.MongoClientSettings
import com.mongodb.ServerApi
import com.mongodb.ServerApiVersion
import com.mongodb.reactivestreams.client.MongoClient
import com.mongodb.reactivestreams.client.MongoClients
import com.wittano.komputer.config.Config
import com.wittano.komputer.config.ConfigLoader
import dagger.Module
import dagger.Provides
import org.bson.Document
import reactor.core.publisher.Mono
import java.time.Duration
import javax.inject.Singleton

@Module
class MongoDbModule {

    @Provides
    @Singleton
    fun client(): MongoClient {
        val config = ConfigLoader.load()

        val serverApi = ServerApi.builder()
            .version(ServerApiVersion.V1)
            .build()

        val settings = MongoClientSettings.builder()
            .applyConnectionString(ConnectionString(config.mongoDbUri))
            .serverApi(serverApi)
            .build()

        val client = MongoClients.create(settings)

        checkDatabaseConnection(client, config)

        return client
    }

    private fun checkDatabaseConnection(client: MongoClient, config: Config) {
        Mono.from(client.getDatabase(config.mongoDbName).runCommand(Document("ping", 1)))
            .timeout(Duration.ofSeconds(2))
            .doOnError {
                throw it
            }
            .block()
    }

}
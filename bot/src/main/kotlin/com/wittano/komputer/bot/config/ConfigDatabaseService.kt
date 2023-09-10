package com.wittano.komputer.bot.config

import com.mongodb.BasicDBObject
import com.wittano.komputer.bot.joke.mongodb.toMonoVoid
import com.wittano.komputer.bot.utils.mongodb.MongoCollectionManager
import com.wittano.komputer.bot.utils.mongodb.database
import com.wittano.komputer.bot.utils.mongodb.getModelCollection
import org.bson.BsonDocument
import org.bson.BsonString
import reactor.core.publisher.Mono
import reactor.kotlin.core.publisher.switchIfEmpty
import reactor.kotlin.core.publisher.toMono

private const val CONFIG_COLLECTION_NAME = "config"

class ConfigDatabaseService {

    private val configCollection = getModelCollection<ServerConfigModel>(CONFIG_COLLECTION_NAME)

    init {
        database.flatMapMany {
            MongoCollectionManager.createCollectionIfDontExist(it, "config")
        }.subscribe()
    }

    operator fun get(guid: String): Mono<ServerConfig> {
        return configCollection.flatMap { collection ->
            val configFilter = getConfigFilter(guid)

            collection.find(configFilter)
                .toMono()
                .switchIfEmpty {
                    val config = ServerConfig().toModel(guid)

                    collection.insertOne(config).toMono().map { config }
                }
                .map {
                    it.toServerConfig()
                }
        }
    }

    fun update(guid: String, config: ServerConfig): Mono<Void> {
        return configCollection.flatMap { collection ->
            val configFilter = BasicDBObject().apply {
                this["guid"] = guid
            }

            Mono.from(collection.find(configFilter))
                .flatMap {
                    val diff = it.createDiff(config)

                    collection.updateOne(configFilter, diff).toMono().toMonoVoid()
                }
                .switchIfEmpty(collection.insertOne(config.toModel(guid)).toMono().toMonoVoid())
                .toMonoVoid()
        }
    }

    private fun getConfigFilter(guid: String) = BasicDBObject().apply {
        this["guid"] = guid
    }

}

private fun ServerConfigModel.createDiff(update: ServerConfig): BsonDocument {
    val document = BsonDocument()

    this.language.takeIf { it != update.language.language }?.also {
        document["language"] = BsonString(it)
    }

    return document
}
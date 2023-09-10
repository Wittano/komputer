package com.wittano.komputer.bot.config

import com.mongodb.BasicDBObject
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

    fun update(guid: String, config: ServerConfig): Mono<ServerConfig> {
        return configCollection.flatMap { collection ->
            val configFilter = getConfigFilter(guid)

            collection.find(configFilter)
                .toMono()
                .flatMap { model ->
                    val diff = model.createDiff(config)

                    collection.updateOne(configFilter, diff).toMono().map { config }
                }
                .switchIfEmpty(collection.insertOne(config.toModel(guid)).toMono().map { config })
        }
    }

    private fun getConfigFilter(guid: String) = BasicDBObject().apply {
        this["guid"] = guid
    }

}

private fun ServerConfigModel.createDiff(update: ServerConfig): BsonDocument {
    val document = BsonDocument()

    this.language.takeIf { it != update.language.language }?.also {
        document["\$set"] = BsonDocument("language", BsonString(it))
    }

    update.roleId?.takeIf { it != this.roleId }?.also {
        document["\$set"] = BsonDocument("role", BsonString(it))
    }

    return document
}
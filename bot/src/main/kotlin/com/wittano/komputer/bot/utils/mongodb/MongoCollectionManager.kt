package com.wittano.komputer.bot.utils.mongodb

import com.mongodb.reactivestreams.client.MongoDatabase
import org.slf4j.LoggerFactory
import reactor.core.publisher.Flux

class MongoCollectionManager private constructor() {
    companion object {
        private val log = LoggerFactory.getLogger(MongoCollectionManager::class.qualifiedName)

        fun createCollectionIfDontExist(database: MongoDatabase, vararg names: String): Flux<Void> {
            val collectionsNames = Flux.from(database.listCollectionNames())
                .collectList()
                .filter {
                    !it.containsAll(names.toList())
                }.flatMapIterable {
                    it
                }

            return collectionsNames
                .flatMap {
                    Flux.from(database.createCollection(it))
                }.doOnError {
                    log.error("Failed to create collection", it)
                }
        }
    }

}
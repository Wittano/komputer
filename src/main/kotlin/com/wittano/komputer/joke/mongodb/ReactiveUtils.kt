package com.wittano.komputer.joke.mongodb

import reactor.core.publisher.Mono

internal fun Mono<*>.toMonoVoid(): Mono<Void> = this.flatMap { Mono.defer { this.then() } }
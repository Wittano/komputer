package com.wittano.komputer.core.config.dagger.key

import dagger.MapKey

@MapKey
@Target(AnnotationTarget.FUNCTION)
@Retention(AnnotationRetention.RUNTIME)
annotation class JokeServiceKey(val value: JokeServiceType)

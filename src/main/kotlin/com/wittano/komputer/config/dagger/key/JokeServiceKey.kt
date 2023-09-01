package com.wittano.komputer.config.dagger.key

import dagger.MapKey

@MapKey
@Target(AnnotationTarget.FUNCTION)
@Retention(AnnotationRetention.RUNTIME)
annotation class JokeServiceKey(val value: JokeServiceType)

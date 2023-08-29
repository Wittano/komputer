package com.wittano.komputer.utils

import java.util.*

fun <T> Optional<T>.toNullable(): T? = if (this.isPresent) {
    this.get()
} else {
    null
}
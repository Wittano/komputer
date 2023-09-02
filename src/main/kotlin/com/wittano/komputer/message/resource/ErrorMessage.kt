package com.wittano.komputer.message.resource

enum class ErrorMessage(val code: String) {
    MISSING_QUESTION_FILED("validation.missing-question"),
    JOKE_NOT_FOUND("joke.not-found"),
    UNSUPPORTED_TYPE("joke.unsupported-type"),
    UNSUPPORTED_CATEGORY("joke.unsupported-category")
}
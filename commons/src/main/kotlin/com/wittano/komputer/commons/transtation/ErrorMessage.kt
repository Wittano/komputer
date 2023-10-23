package com.wittano.komputer.commons.transtation

enum class ErrorMessage(val code: String) {
    MISSING_QUESTION_FILED("validation.missing-question"),
    JOKE_NOT_FOUND("joke.not-found"),
    UNSUPPORTED_TYPE("joke.unsupported-type"),
    UNSUPPORTED_CATEGORY("joke.unsupported-category"),
    JOKE_ID_INVALID("joke.invalid-joke-id"),
    GENERAL_ERROR("general.error"),
    ACCESS_DENIED("general.access-denied"),
    CONFIG_UPDATE_FAILED("config.update-failed"),
}
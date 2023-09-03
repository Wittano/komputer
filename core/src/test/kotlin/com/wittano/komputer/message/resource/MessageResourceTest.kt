package com.wittano.komputer.message.resource

import org.junit.jupiter.api.Test
import java.util.*
import kotlin.test.assertEquals

class MessageResourceTest {

    @Test
    fun getErrorMessage() {
        val message = MessageResource.get(ErrorMessage.MISSING_QUESTION_FILED, Locale.ENGLISH)

        assertEquals("TEST", message)
    }

}
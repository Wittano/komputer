package com.wittano.komputer.message.resource

import com.wittano.komputer.core.message.resource.ErrorMessage
import com.wittano.komputer.core.message.resource.MessageResource
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test
import java.util.*

class MessageResourceTest {

    @Test
    fun getErrorMessage() {
        val message = MessageResource.get(ErrorMessage.MISSING_QUESTION_FILED, Locale.ENGLISH)

        assertEquals("TEST", message)
    }

}
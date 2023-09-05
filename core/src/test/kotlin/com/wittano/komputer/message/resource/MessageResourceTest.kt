package com.wittano.komputer.message.resource

import com.wittano.komputer.core.message.resource.ErrorMessage
import com.wittano.komputer.core.message.resource.getErrorMessage
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test
import java.util.*

class MessageResourceTest {

    @Test
    fun getErrorMessage_success() {
        val message = getErrorMessage(ErrorMessage.MISSING_QUESTION_FILED, Locale.ENGLISH)

        assertEquals("TEST", message)
    }

}
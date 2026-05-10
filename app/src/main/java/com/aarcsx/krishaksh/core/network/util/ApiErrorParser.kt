package com.aarcsx.krishaksh.core.network.util

import com.aarcsx.krishaksh.core.network.dto.ErrorEnvelope
import com.google.gson.Gson
import retrofit2.Response

object ApiErrorParser {
    private val gson = Gson()

    fun parse(response: Response<*>): String {
        return try {
            val raw = response.errorBody()?.string().orEmpty()
            if (raw.isBlank()) "Unknown API error"
            else gson.fromJson(raw, ErrorEnvelope::class.java)?.message ?: "Unknown API error"
        } catch (_: Exception) {
            "Unknown API error"
        }
    }
}

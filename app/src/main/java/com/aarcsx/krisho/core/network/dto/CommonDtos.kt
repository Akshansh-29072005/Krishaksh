package com.aarcsx.krisho.core.network.dto

data class ApiEnvelope<T>(val success: Boolean, val message: String, val data: T?)
data class ErrorEnvelope(val success: Boolean? = null, val message: String? = null)

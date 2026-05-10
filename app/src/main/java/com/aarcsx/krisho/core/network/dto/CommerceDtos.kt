package com.aarcsx.krisho.core.network.dto

data class CreateOrderRequestDto(val shipping_metadata: Map<String, Any> = emptyMap(), val notes: String? = null)
data class OrderDto(val id: String, val status: String, val grand_total: Double, val currency: String, val created_at: String)
data class CreatePaymentOrderRequestDto(val order_id: String)
data class RazorpayOrderDto(val key_id: String, val order_id: String, val amount: Long, val currency: String, val receipt: String)

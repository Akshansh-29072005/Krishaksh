package com.aarcsx.krisho.core.network.api

import com.aarcsx.krisho.core.network.dto.*
import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.POST
import retrofit2.http.Path

interface CommerceApiService {
    @POST("orders")
    suspend fun createOrder(@Body body: CreateOrderRequestDto): Response<ApiEnvelope<OrderDto>>

    @GET("orders/{id}")
    suspend fun getOrder(@Path("id") id: String): Response<ApiEnvelope<OrderDto>>

    @POST("payments/create-order")
    suspend fun createPaymentOrder(@Body body: CreatePaymentOrderRequestDto): Response<ApiEnvelope<RazorpayOrderDto>>
}

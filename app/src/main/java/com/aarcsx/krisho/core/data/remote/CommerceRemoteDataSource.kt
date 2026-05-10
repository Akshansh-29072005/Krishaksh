package com.aarcsx.krisho.core.data.remote

import com.aarcsx.krisho.core.common.ApiResult
import com.aarcsx.krisho.core.network.api.CommerceApiService
import com.aarcsx.krisho.core.network.dto.*
import javax.inject.Inject

class CommerceRemoteDataSource @Inject constructor(
    private val api: CommerceApiService
) : BaseRemoteDataSource() {
    suspend fun createOrder(notes: String?): ApiResult<ApiEnvelope<OrderDto>> = call { api.createOrder(CreateOrderRequestDto(notes = notes)) }
    suspend fun getOrder(orderId: String): ApiResult<ApiEnvelope<OrderDto>> = call { api.getOrder(orderId) }
    suspend fun createPaymentOrder(orderId: String): ApiResult<ApiEnvelope<RazorpayOrderDto>> = call { api.createPaymentOrder(CreatePaymentOrderRequestDto(orderId)) }
}

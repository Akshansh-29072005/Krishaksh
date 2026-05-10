package com.aarcsx.krishaksh.core.data.remote

import com.aarcsx.krishaksh.core.common.ApiResult
import com.aarcsx.krishaksh.core.network.api.CommerceApiService
import com.aarcsx.krishaksh.core.network.dto.*
import javax.inject.Inject

class CommerceRemoteDataSource @Inject constructor(
    private val api: CommerceApiService
) : BaseRemoteDataSource() {
    suspend fun createOrder(notes: String?): ApiResult<ApiEnvelope<OrderDto>> = call { api.createOrder(CreateOrderRequestDto(notes = notes)) }
    suspend fun getOrder(orderId: String): ApiResult<ApiEnvelope<OrderDto>> = call { api.getOrder(orderId) }
    suspend fun createPaymentOrder(orderId: String): ApiResult<ApiEnvelope<RazorpayOrderDto>> = call { api.createPaymentOrder(CreatePaymentOrderRequestDto(orderId)) }
}

package com.aarcsx.krisho.core.payments

import com.aarcsx.krisho.core.common.ApiResult
import com.aarcsx.krisho.core.data.remote.CommerceRemoteDataSource
import com.aarcsx.krisho.core.network.dto.OrderDto
import com.aarcsx.krisho.core.network.dto.RazorpayOrderDto
import kotlinx.coroutines.delay
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class PaymentRepository @Inject constructor(
    private val commerceRemote: CommerceRemoteDataSource
) {
    suspend fun createOrder(notes: String? = null): ApiResult<OrderDto> {
        return when (val res = commerceRemote.createOrder(notes)) {
            is ApiResult.Success -> res.data.data?.let { ApiResult.Success(it) } ?: ApiResult.Error(message = "Invalid order response")
            is ApiResult.Error -> res
            ApiResult.Loading -> ApiResult.Loading
        }
    }

    suspend fun createRazorpayOrder(orderId: String): ApiResult<RazorpayOrderDto> {
        return when (val res = commerceRemote.createPaymentOrder(orderId)) {
            is ApiResult.Success -> res.data.data?.let { ApiResult.Success(it) } ?: ApiResult.Error(message = "Invalid payment order response")
            is ApiResult.Error -> res
            ApiResult.Loading -> ApiResult.Loading
        }
    }

    suspend fun pollOrderState(orderId: String, retries: Int = 8, delayMs: Long = 1500): ApiResult<OrderDto> {
        repeat(retries) {
            when (val state = commerceRemote.getOrder(orderId)) {
                is ApiResult.Success -> {
                    val ord = state.data.data ?: return ApiResult.Error(message = "Invalid order state")
                    if (ord.status in setOf("paid", "failed", "cancelled", "refunded")) return ApiResult.Success(ord)
                }
                is ApiResult.Error -> if (it == retries - 1) return state
                ApiResult.Loading -> Unit
            }
            delay(delayMs)
        }
        return ApiResult.Error(message = "Order state polling timeout")
    }
}

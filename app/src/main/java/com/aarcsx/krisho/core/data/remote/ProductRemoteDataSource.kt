package com.aarcsx.krisho.core.data.remote

import com.aarcsx.krisho.core.common.ApiResult
import com.aarcsx.krisho.core.network.api.ProductApiService
import com.aarcsx.krisho.core.network.dto.*
import javax.inject.Inject

class ProductRemoteDataSource @Inject constructor(
    private val api: ProductApiService
) : BaseRemoteDataSource() {
    suspend fun getProducts(): ApiResult<ApiEnvelope<List<ProductDto>>> = call { api.getProducts() }
    suspend fun getProduct(id: String): ApiResult<ApiEnvelope<ProductDto>> = call { api.getProduct(id) }
}

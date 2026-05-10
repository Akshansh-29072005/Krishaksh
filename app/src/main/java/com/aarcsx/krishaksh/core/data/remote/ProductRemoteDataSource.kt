package com.aarcsx.krishaksh.core.data.remote

import com.aarcsx.krishaksh.core.common.ApiResult
import com.aarcsx.krishaksh.core.network.api.ProductApiService
import com.aarcsx.krishaksh.core.network.dto.*
import javax.inject.Inject

class ProductRemoteDataSource @Inject constructor(
    private val api: ProductApiService
) : BaseRemoteDataSource() {
    suspend fun getProducts(): ApiResult<ApiEnvelope<List<ProductDto>>> = call { api.getProducts() }
    suspend fun getProduct(id: String): ApiResult<ApiEnvelope<ProductDto>> = call { api.getProduct(id) }
}

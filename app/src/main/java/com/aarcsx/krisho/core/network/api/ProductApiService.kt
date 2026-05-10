package com.aarcsx.krisho.core.network.api

import com.aarcsx.krisho.core.network.dto.*
import retrofit2.Response
import retrofit2.http.GET
import retrofit2.http.Path

interface ProductApiService {
    @GET("products")
    suspend fun getProducts(): Response<ApiEnvelope<List<ProductDto>>>

    @GET("products/{id}")
    suspend fun getProduct(@Path("id") id: String): Response<ApiEnvelope<ProductDto>>
}

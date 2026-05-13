package com.aarcsx.krisho.core.network.api

import com.aarcsx.krisho.core.network.dto.AppConfigDto
import com.aarcsx.krisho.core.network.dto.CropDto
import com.aarcsx.krisho.core.network.dto.ApiEnvelope
import retrofit2.Response
import retrofit2.http.GET

interface AppApiService {
    @GET("crops")
    suspend fun getCrops(): Response<ApiEnvelope<List<CropDto>>>

    @GET("app-config")
    suspend fun getAppConfig(): Response<ApiEnvelope<AppConfigDto>>
}

package com.aarcsx.krisho.core.network.api

import com.aarcsx.krisho.core.network.dto.*
import okhttp3.RequestBody
import retrofit2.Response
import retrofit2.http.*

interface ScanApiService {
    @GET("scans/upload-url")
    suspend fun getUploadUrl(@Query("content_type") contentType: String): Response<ApiEnvelope<PresignedUploadDto>>

    @POST("scans")
    suspend fun createScan(@Body body: CreateScanRequestDto): Response<ApiEnvelope<ScanDto>>

    @GET("scans/{id}")
    suspend fun getScan(@Path("id") id: String): Response<ApiEnvelope<ScanDto>>

    @GET("recommendations/{scanId}")
    suspend fun getRecommendations(@Path("scanId") scanId: String): Response<ApiEnvelope<RecommendationDto>>

    @PUT
    suspend fun uploadToS3(
        @Url uploadUrl: String,
        @Header("Content-Type") contentType: String,
        @Body body: RequestBody
    ): Response<Unit>
}

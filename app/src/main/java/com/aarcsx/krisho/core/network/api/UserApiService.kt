package com.aarcsx.krisho.core.network.api

import com.aarcsx.krisho.core.network.dto.*
import okhttp3.MultipartBody
import retrofit2.Response
import retrofit2.http.*

interface UserApiService {
    @GET("users/me")
    suspend fun me(): Response<ApiEnvelope<UserMeDto>>

    @PUT("users/me/phone")
    suspend fun updatePhone(@Body body: UpdatePhoneDto): Response<ApiEnvelope<UserMeDto>>

    @GET("support/tickets")
    suspend fun getTickets(): Response<ApiEnvelope<List<SupportTicketDto>>>

    @POST("support/tickets")
    suspend fun createTicket(@Body body: CreateTicketDto): Response<ApiEnvelope<SupportTicketDto>>

    @GET("support/tickets/{id}")
    suspend fun getTicket(@Path("id") id: String): Response<ApiEnvelope<SupportThreadDto>>

    @POST("support/tickets/{id}/messages")
    suspend fun sendMessage(@Path("id") id: String, @Body body: SendMessageDto): Response<ApiEnvelope<SupportMessageDto>>

    @POST("support/tickets/{id}/callback")
    suspend fun requestCallback(@Path("id") id: String): Response<ApiEnvelope<SupportTicketDto>>

    @Multipart
    @POST("support/tickets/{id}/voice")
    suspend fun uploadVoice(@Path("id") id: String, @Part voice: MultipartBody.Part): Response<ApiEnvelope<Any>>

    @GET("ads/featured")
    suspend fun featuredAds(): Response<ApiEnvelope<List<AdCampaignDto>>>

    @PUT("users/me/crops")
    suspend fun updateCrops(@Body body: UpdateCropsDto): Response<ApiEnvelope<UserMeDto>>
}

data class UpdatePhoneDto(
    val phone: String
)

data class UpdateCropsDto(
    val crops: List<String>
)

package com.aarcsx.krishaksh.core.network.api

import com.aarcsx.krishaksh.core.network.dto.*
import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.POST

interface AuthApiService {
    @POST("auth/google")
    suspend fun googleLogin(@Body body: GoogleLoginRequest): Response<ApiEnvelope<AuthTokensDto>>

    @POST("auth/refresh")
    suspend fun refresh(@Body body: RefreshTokenRequest): Response<ApiEnvelope<AuthTokensDto>>

    @POST("auth/logout")
    suspend fun logout(@Body body: LogoutRequest): Response<ApiEnvelope<Map<String, Any>>>
}

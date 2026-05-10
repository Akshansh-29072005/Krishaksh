package com.aarcsx.krishaksh.core.data.remote

import com.aarcsx.krishaksh.core.common.ApiResult
import com.aarcsx.krishaksh.core.network.api.AuthApiService
import com.aarcsx.krishaksh.core.network.dto.*
import javax.inject.Inject

class AuthRemoteDataSource @Inject constructor(
    private val api: AuthApiService
) : BaseRemoteDataSource() {
    suspend fun googleLogin(idToken: String): ApiResult<ApiEnvelope<AuthTokensDto>> = call { api.googleLogin(GoogleLoginRequest(idToken)) }
    suspend fun refreshToken(refreshToken: String): ApiResult<ApiEnvelope<AuthTokensDto>> = call { api.refresh(RefreshTokenRequest(refreshToken)) }
    suspend fun logout(refreshToken: String): ApiResult<ApiEnvelope<Map<String, Any>>> = call { api.logout(LogoutRequest(refreshToken)) }
}

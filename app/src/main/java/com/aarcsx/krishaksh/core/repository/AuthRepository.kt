package com.aarcsx.krishaksh.core.repository

import com.aarcsx.krishaksh.core.common.ApiResult
import com.aarcsx.krishaksh.core.data.remote.AuthRemoteDataSource
import com.aarcsx.krishaksh.core.local.datastore.PreferencesManager
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.first
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class AuthRepository @Inject constructor(
    private val preferencesManager: PreferencesManager,
    private val remote: AuthRemoteDataSource
) {
    val jwtToken: Flow<String?> = preferencesManager.jwtToken
    val refreshToken: Flow<String?> = preferencesManager.refreshToken
    val languageSetting: Flow<String> = preferencesManager.languageSetting

    suspend fun setLanguage(langCode: String) = preferencesManager.saveLanguage(langCode)

    suspend fun loginWithGoogle(idToken: String): ApiResult<Unit> {
        return when (val result = remote.googleLogin(idToken)) {
            is ApiResult.Success -> {
                val tokens = result.data.data ?: return ApiResult.Error(message = "Invalid auth response")
                preferencesManager.saveTokens(tokens.access_token, tokens.refresh_token)
                ApiResult.Success(Unit)
            }
            is ApiResult.Error -> result
            ApiResult.Loading -> ApiResult.Loading
        }
    }

    suspend fun logout(): ApiResult<Unit> {
        val refresh = preferencesManager.refreshToken.first()
        if (!refresh.isNullOrBlank()) remote.logout(refresh)
        preferencesManager.clearTokens()
        return ApiResult.Success(Unit)
    }
}

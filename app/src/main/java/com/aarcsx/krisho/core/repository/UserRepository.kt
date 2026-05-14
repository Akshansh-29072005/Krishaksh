package com.aarcsx.krisho.core.repository

import com.aarcsx.krisho.core.common.ApiResult
import com.aarcsx.krisho.core.local.datastore.PreferencesManager
import com.aarcsx.krisho.core.network.api.UpdateCropsDto
import com.aarcsx.krisho.core.network.api.UpdatePhoneDto
import com.aarcsx.krisho.core.network.api.UserApiService
import com.aarcsx.krisho.core.network.dto.UserMeDto
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.flow
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class UserRepository @Inject constructor(
    private val apiService: UserApiService,
    private val preferencesManager: PreferencesManager
) {
    fun getProfile(): Flow<ApiResult<UserMeDto>> = flow {
        emit(ApiResult.Loading)
        try {
            val response = apiService.me()
            if (response.isSuccessful) {
                response.body()?.data?.let {
                    emit(ApiResult.Success(it))
                } ?: emit(ApiResult.Error(message = "Empty response"))
            } else {
                emit(ApiResult.Error(message = response.message()))
            }
        } catch (e: Exception) {
            emit(ApiResult.Error(message = e.localizedMessage ?: "Unknown error"))
        }
    }

    suspend fun saveLanguage(lang: String) {
        preferencesManager.saveLanguage(lang)
    }

    suspend fun updatePhone(phone: String): ApiResult<UserMeDto> {
        return try {
            val response = apiService.updatePhone(UpdatePhoneDto(phone))
            if (response.isSuccessful) {
                response.body()?.data?.let { ApiResult.Success(it) }
                    ?: ApiResult.Error(message = "Empty response")
            } else {
                ApiResult.Error(message = response.message())
            }
        } catch (e: Exception) {
            ApiResult.Error(message = e.localizedMessage ?: "Unknown error")
        }
    }

    suspend fun updateCrops(crops: List<String>): ApiResult<UserMeDto> {
        return try {
            val response = apiService.updateCrops(UpdateCropsDto(crops))
            if (response.isSuccessful) {
                response.body()?.data?.let { ApiResult.Success(it) }
                    ?: ApiResult.Error(message = "Empty response")
            } else {
                ApiResult.Error(message = response.message())
            }
        } catch (e: Exception) {
            ApiResult.Error(message = e.localizedMessage ?: "Unknown error")
        }
    }

    val languageSetting: Flow<String> = preferencesManager.languageSetting

    suspend fun logout() {
        preferencesManager.clearTokens()
    }
}

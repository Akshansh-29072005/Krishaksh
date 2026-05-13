package com.aarcsx.krisho.core.repository

import com.aarcsx.krisho.core.common.ApiResult
import com.aarcsx.krisho.core.local.datastore.PreferencesManager
import com.aarcsx.krisho.core.network.api.AppApiService
import com.aarcsx.krisho.core.network.dto.AppConfigDto
import com.aarcsx.krisho.core.network.dto.CropDto
import kotlinx.coroutines.flow.first
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class AppRepository @Inject constructor(
    private val apiService: AppApiService,
    private val preferencesManager: PreferencesManager
) {
    suspend fun getCrops(): ApiResult<List<CropDto>> {
        return try {
            val response = apiService.getCrops()
            if (response.isSuccessful) {
                response.body()?.data?.let { ApiResult.Success(it) } ?: ApiResult.Error(message = "Empty crop list")
            } else {
                ApiResult.Error(message = response.message())
            }
        } catch (e: Exception) {
            ApiResult.Error(message = e.localizedMessage ?: "Unknown error")
        }
    }

    suspend fun getAppConfig(): ApiResult<AppConfigDto> {
        return try {
            val response = apiService.getAppConfig()
            if (response.isSuccessful) {
                response.body()?.data?.let {
                    preferencesManager.saveAppConfig(it)
                    ApiResult.Success(it)
                } ?: loadCachedAppConfigOrError("Empty app config")
            } else {
                loadCachedAppConfigOrError("Unable to load app config: ${response.message()}")
            }
        } catch (e: Exception) {
            loadCachedAppConfigOrError(e.localizedMessage ?: "Unknown error")
        }
    }

    private suspend fun loadCachedAppConfigOrError(message: String): ApiResult<AppConfigDto> {
        return preferencesManager.cachedAppConfig.first()?.let { ApiResult.Success(it) }
            ?: ApiResult.Error(message = message)
    }
}

package com.aarcsx.krisho.core.data.remote

import com.aarcsx.krisho.core.common.ApiResult
import com.aarcsx.krisho.core.network.api.ScanApiService
import com.aarcsx.krisho.core.network.dto.*
import okhttp3.MediaType.Companion.toMediaTypeOrNull
import okhttp3.RequestBody
import okhttp3.RequestBody.Companion.toRequestBody
import javax.inject.Inject

class ScanRemoteDataSource @Inject constructor(
    private val api: ScanApiService
) : BaseRemoteDataSource() {
    suspend fun getPresignedUrl(contentType: String): ApiResult<ApiEnvelope<PresignedUploadDto>> = call { api.getUploadUrl(contentType) }
    suspend fun uploadToS3(uploadUrl: String, bytes: ByteArray, contentType: String): ApiResult<Unit> {
        return when (val res = call {
            api.uploadToS3(
                uploadUrl,
                contentType,
                bytes.toRequestBody(contentType.toMediaTypeOrNull())
            )
        }) {
            is ApiResult.Success -> ApiResult.Success(Unit)
            is ApiResult.Error -> res
            ApiResult.Loading -> ApiResult.Loading
        }
    }
    suspend fun createScan(cropType: String, imageKey: String): ApiResult<ApiEnvelope<ScanDto>> = call { api.createScan(CreateScanRequestDto(cropType, imageKey)) }
    suspend fun getScan(scanId: String): ApiResult<ApiEnvelope<ScanDto>> = call { api.getScan(scanId) }
    suspend fun getRecommendations(scanId: String): ApiResult<ApiEnvelope<RecommendationDto>> = call { api.getRecommendations(scanId) }
    suspend fun getScanHistory(): ApiResult<ApiEnvelope<List<ScanDto>>> = call { api.getScanHistory() }
}

package com.aarcsx.krisho.core.repository

import com.aarcsx.krisho.core.common.ApiResult
import com.aarcsx.krisho.core.data.remote.ScanRemoteDataSource
import com.aarcsx.krisho.core.local.room.dao.ScanDao
import com.aarcsx.krisho.core.local.room.entity.ScanEntity
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.Flow
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class ScanRepository @Inject constructor(
    private val scanDao: ScanDao,
    private val remote: ScanRemoteDataSource,
    private val offlineSyncRepository: OfflineSyncRepository
) {
    fun getRecentScans(): Flow<List<ScanEntity>> = scanDao.getAllScans()

    suspend fun runScanLifecycle(cropName: String, imageBytes: ByteArray, contentType: String = "image/jpeg"): ApiResult<ScanEntity> {
        return try {
            val uploadRes = remote.getPresignedUrl(contentType)
            if (uploadRes is ApiResult.Error) return uploadRes
            
            val upload = when (uploadRes) {
                is ApiResult.Success -> uploadRes.data.data ?: return ApiResult.Error(message = "Failed upload-url response")
                else -> return ApiResult.Error(message = "Unexpected state")
            }

            when (val uploadResult = remote.uploadToS3(upload.presigned_url, imageBytes, contentType)) {
                is ApiResult.Error -> return uploadResult
                else -> Unit
            }

            val scan = when (val createRes = remote.createScan(cropName, upload.image_key)) {
                is ApiResult.Success -> createRes.data.data ?: return ApiResult.Error(message = "Failed create-scan response")
                is ApiResult.Error -> return createRes
                ApiResult.Loading -> return ApiResult.Error(message = "Unexpected loading state")
            }

            // Increase polling window for slower AI providers
            val final = pollScan(scan.id) ?: scan
            
            val rec = when (val recRes = remote.getRecommendations(scan.id)) {
                is ApiResult.Success -> recRes.data.data
                else -> null
            }

            val entity = ScanEntity(
                cropName = final.crop_type,
                imageUrl = final.image_url,
                remoteScanId = final.id,
                predictionStatus = final.prediction_status,
                diseaseName = rec?.disease?.name ?: final.disease_name ?: "Pending",
                confidence = final.confidence_score?.let { "${"%.1f".format(it * 100)}%" } ?: "-",
                symptoms = rec?.disease?.symptoms ?: final.ai_symptoms?.joinToString("\n• ") ?: "",
                prevention = rec?.disease?.prevention.orEmpty(),
                treatment = rec?.disease?.treatment.orEmpty(),
                recommendationTitle = rec?.recommended_products?.firstOrNull()?.name ?: "",
                recommendationDesc = rec?.recommended_products?.firstOrNull()?.description.orEmpty(),
                capturedAt = System.currentTimeMillis()
            )
            scanDao.insertScan(entity)
            ApiResult.Success(entity)
        } catch (e: Exception) {
            e.printStackTrace()
            ApiResult.Error(message = "Scan Lifecycle Error: ${e.message}")
        }
    }

    private suspend fun pollScan(scanId: String, retries: Int = 15, delayMs: Long = 3000): com.aarcsx.krisho.core.network.dto.ScanDto? {
        repeat(retries) {
            when (val statusRes = remote.getScan(scanId)) {
                is ApiResult.Success -> {
                    val scan = statusRes.data.data ?: return null
                    if (scan.prediction_status.equals("COMPLETED", true) ||
                        scan.prediction_status.equals("FAILED", true)) {
                        return scan
                    }
                }
                else -> Unit
            }
            delay(delayMs)
        }
        return null
    }
}

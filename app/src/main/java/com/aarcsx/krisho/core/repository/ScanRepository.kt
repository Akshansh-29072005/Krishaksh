package com.aarcsx.krisho.core.repository

import android.content.Context
import com.aarcsx.krisho.core.common.ApiResult
import com.aarcsx.krisho.core.data.remote.ScanRemoteDataSource
import com.aarcsx.krisho.core.local.room.dao.ScanDao
import com.aarcsx.krisho.core.local.room.entity.ScanEntity
import com.aarcsx.krisho.core.util.NetworkMonitor
import com.aarcsx.krisho.core.work.OfflineScanWorker
import com.google.gson.Gson
import dagger.hilt.android.qualifiers.ApplicationContext
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.Flow
import java.io.File
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class ScanRepository @Inject constructor(
    private val scanDao: ScanDao,
    private val remote: ScanRemoteDataSource,
    private val offlineSyncRepository: OfflineSyncRepository,
    @ApplicationContext private val context: Context,
    private val networkMonitor: NetworkMonitor
) {

    fun getRecentScans(): Flow<List<ScanEntity>> = scanDao.getAllScans()

    suspend fun refreshScans() {
        when (val res = remote.getScanHistory()) {
            is ApiResult.Success -> {
                val pending = scanDao.getPendingLocalScans()
                val pendingOffline = pending.filter { it.remoteScanId.isNullOrBlank() }
                val pendingRemote = pending.filter { !it.remoteScanId.isNullOrBlank() }
                    .associateBy { it.remoteScanId!! }

                val scans = res.data.data.orEmpty()
                val entities = scans.map { dto ->
                    ScanEntity(
                        cropName = dto.crop_type,
                        imageUrl = dto.image_url,
                        remoteScanId = dto.id,
                        predictionStatus = dto.prediction_status,
                        diseaseName = dto.disease_name ?: "Pending",
                        confidence = dto.confidence_score?.let { "${"%.1f".format(it * 100)}%" } ?: "-",
                        symptoms = dto.ai_symptoms?.joinToString("\n• ") ?: "",
                        prevention = "",
                        treatment = "",
                        recommendationTitle = "",
                        recommendationDesc = "",
                        capturedAt = System.currentTimeMillis()
                    )
                }

                val mergedScans = entities.map { serverScan ->
                    pendingRemote[serverScan.remoteScanId] ?: serverScan
                } + pendingOffline

                scanDao.clearHistory()
                scanDao.insertScans(mergedScans)
            }
            else -> Unit
        }
    }

    suspend fun runScanLifecycle(cropName: String, imageBytes: ByteArray, contentType: String = "image/jpeg"): ApiResult<ScanEntity> {
        return try {
            if (!networkMonitor.isConnected()) {
                return queueOfflineScan(cropName, imageBytes, contentType)
            }

            val uploadRes = remote.getPresignedUrl(contentType)
            if (uploadRes is ApiResult.Error) {
                return queueOfflineScan(cropName, imageBytes, contentType)
            }

            val upload = when (uploadRes) {
                is ApiResult.Success -> uploadRes.data.data ?: return queueOfflineScan(cropName, imageBytes, contentType)
                else -> return queueOfflineScan(cropName, imageBytes, contentType)
            }

            when (val uploadResult = remote.uploadToS3(upload.presigned_url, imageBytes, contentType)) {
                is ApiResult.Error -> return queueOfflineScan(cropName, imageBytes, contentType)
                else -> Unit
            }

            val scan = when (val createRes = remote.createScan(cropName, upload.image_key)) {
                is ApiResult.Success -> createRes.data.data ?: return queueOfflineScan(cropName, imageBytes, contentType)
                is ApiResult.Error -> return queueOfflineScan(cropName, imageBytes, contentType)
                ApiResult.Loading -> return queueOfflineScan(cropName, imageBytes, contentType)
            }

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
            queueOfflineScan(cropName, imageBytes, contentType)
        }
    }

    private suspend fun queueOfflineScan(cropName: String, imageBytes: ByteArray, contentType: String): ApiResult<ScanEntity> {
        return try {
            val imageFile = createOfflineImageFile(imageBytes)
            val scan = ScanEntity(
                cropName = cropName,
                imageUrl = imageFile.toURI().toString(),
                predictionStatus = "QUEUED",
                diseaseName = "Pending",
                confidence = "-",
                symptoms = "",
                prevention = "",
                treatment = "",
                recommendationTitle = "",
                recommendationDesc = "",
                capturedAt = System.currentTimeMillis()
            )
            val insertedId = scanDao.insertScan(scan).toInt()
            val queuedScan = scan.copy(id = insertedId)
            val payload = OfflineScanPayload(cropName, imageFile.absolutePath, contentType, insertedId.toLong())
            offlineSyncRepository.enqueue("SCAN_UPLOAD", Gson().toJson(payload))
            OfflineScanWorker.enqueue(context)
            ApiResult.Success(queuedScan)
        } catch (e: Exception) {
            e.printStackTrace()
            ApiResult.Error(message = "Unable to queue scan offline: ${e.message}")
        }
    }

    private fun createOfflineImageFile(imageBytes: ByteArray): File {
        val directory = File(context.filesDir, "offline_scans")
        if (!directory.exists()) {
            directory.mkdirs()
        }
        val file = File(directory, "scan_${System.currentTimeMillis()}.jpg")
        file.writeBytes(imageBytes)
        return file
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

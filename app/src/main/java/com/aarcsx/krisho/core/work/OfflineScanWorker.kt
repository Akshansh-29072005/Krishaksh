package com.aarcsx.krisho.core.work

import android.content.Context
import androidx.room.Room
import androidx.work.*
import com.aarcsx.krisho.BuildConfig
import com.aarcsx.krisho.core.auth.TokenManager
import com.aarcsx.krisho.core.data.remote.ScanRemoteDataSource
import com.aarcsx.krisho.core.local.datastore.PreferencesManager
import com.aarcsx.krisho.core.local.room.KrishoDatabase
import com.aarcsx.krisho.core.local.room.entity.ScanEntity
import com.aarcsx.krisho.core.repository.OfflineScanPayload
import com.google.gson.Gson
import kotlinx.coroutines.runBlocking
import okhttp3.Interceptor
import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import java.io.File
import java.util.concurrent.TimeUnit

class OfflineScanWorker(
    context: Context,
    params: WorkerParameters
) : CoroutineWorker(context, params) {

    override suspend fun doWork(): Result {
        val db = Room.databaseBuilder(applicationContext, KrishoDatabase::class.java, "krisho_db")
            .fallbackToDestructiveMigration()
            .build()

        val offlineDao = db.offlineSyncDao()
        val scanDao = db.scanDao()
        val pendingTasks = offlineDao.getAllPendingOnce()

        if (pendingTasks.isEmpty()) {
            return Result.success()
        }

        val api = buildScanApiService()
        val remote = ScanRemoteDataSource(api)
        val gson = Gson()

        pendingTasks.forEach { task ->
            when (task.type) {
                "SCAN_UPLOAD" -> {
                    val payload = try {
                        gson.fromJson(task.payload, OfflineScanPayload::class.java)
                    } catch (e: Exception) {
                        offlineDao.deleteById(task.id)
                        return@forEach
                    }

                    val imageFile = File(payload.imageFilePath)
                    if (!imageFile.exists()) {
                        offlineDao.deleteById(task.id)
                        return@forEach
                    }

                    val uploadRes = remote.getPresignedUrl(payload.contentType)
                    if (uploadRes is com.aarcsx.krisho.core.common.ApiResult.Error) {
                        return Result.retry()
                    }

                    val upload = (uploadRes as com.aarcsx.krisho.core.common.ApiResult.Success).data.data
                        ?: return Result.retry()

                    val rawBytes = imageFile.readBytes()
                    when (val uploadResult = remote.uploadToS3(upload.presigned_url, rawBytes, payload.contentType)) {
                        is com.aarcsx.krisho.core.common.ApiResult.Error -> return Result.retry()
                        else -> Unit
                    }

                    val createRes = remote.createScan(payload.cropName, upload.image_key)
                    if (createRes is com.aarcsx.krisho.core.common.ApiResult.Error) {
                        return Result.retry()
                    }
                    val scanDto = (createRes as com.aarcsx.krisho.core.common.ApiResult.Success).data.data
                        ?: return Result.retry()

                    val finalScan = pollScan(api, scanDto.id) ?: scanDto

                    try {
                        val localScan = scanDao.getById(payload.localScanId.toInt())
                        if (localScan != null) {
                            val updated = localScan.copy(
                                remoteScanId = finalScan.id,
                                predictionStatus = finalScan.prediction_status,
                                imageUrl = finalScan.image_url,
                                diseaseName = finalScan.disease_name ?: localScan.diseaseName,
                                confidence = finalScan.confidence_score?.let { "${"%.1f".format(it * 100)}%" } ?: localScan.confidence,
                                symptoms = finalScan.ai_symptoms?.joinToString("\n• ") ?: localScan.symptoms
                            )
                            scanDao.updateScan(updated)
                        }
                    } catch (_: Exception) {
                        // If updating fails, continue cleaning queue and removing local temp file.
                    }

                    offlineDao.deleteById(task.id)
                    if (imageFile.exists()) {
                        imageFile.delete()
                    }
                }
                else -> offlineDao.deleteById(task.id)
            }
        }

        return Result.success()
    }

    private suspend fun pollScan(api: com.aarcsx.krisho.core.network.api.ScanApiService, scanId: String): com.aarcsx.krisho.core.network.dto.ScanDto? {
        repeat(15) {
            when (val scanRes = ScanRemoteDataSource(api).getScan(scanId)) {
                is com.aarcsx.krisho.core.common.ApiResult.Success -> {
                    val scan = scanRes.data.data
                    if (scan != null && (scan.prediction_status.equals("COMPLETED", true) || scan.prediction_status.equals("FAILED", true))) {
                        return scan
                    }
                }
                else -> Unit
            }
            kotlinx.coroutines.delay(3000)
        }
        return null
    }

    private fun buildScanApiService(): com.aarcsx.krisho.core.network.api.ScanApiService {
        val tokenManager = TokenManager(PreferencesManager(applicationContext))
        val logging = HttpLoggingInterceptor().apply {
            level = if (BuildConfig.DEBUG) HttpLoggingInterceptor.Level.BODY else HttpLoggingInterceptor.Level.NONE
        }

        val okHttpClient = OkHttpClient.Builder()
            .addInterceptor { chain ->
                val request = chain.request()
                val builder = request.newBuilder().addHeader("Accept", "application/json")

                if (!request.url.host.contains("amazonaws.com")) {
                    val token = runBlocking { tokenManager.accessToken() }
                    if (!token.isNullOrBlank()) {
                        builder.addHeader("Authorization", "Bearer $token")
                    }
                }

                chain.proceed(builder.build())
            }
            .addInterceptor(logging)
            .connectTimeout(30, TimeUnit.SECONDS)
            .readTimeout(60, TimeUnit.SECONDS)
            .writeTimeout(60, TimeUnit.SECONDS)
            .build()

        return Retrofit.Builder()
            .baseUrl(BuildConfig.BASE_URL)
            .client(okHttpClient)
            .addConverterFactory(GsonConverterFactory.create())
            .build()
            .create(com.aarcsx.krisho.core.network.api.ScanApiService::class.java)
    }

    companion object {
        fun enqueue(context: Context) {
            val workRequest = OneTimeWorkRequestBuilder<OfflineScanWorker>()
                .setConstraints(
                    Constraints.Builder()
                        .setRequiredNetworkType(NetworkType.CONNECTED)
                        .build()
                )
                .setBackoffCriteria(BackoffPolicy.EXPONENTIAL, 30, TimeUnit.SECONDS)
                .build()

            WorkManager.getInstance(context).enqueueUniqueWork(
                "offline_scan_sync",
                ExistingWorkPolicy.KEEP,
                workRequest
            )
        }
    }
}

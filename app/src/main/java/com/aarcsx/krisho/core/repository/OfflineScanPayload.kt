package com.aarcsx.krisho.core.repository

import androidx.annotation.Keep

@Keep
data class OfflineScanPayload(
    val cropName: String,
    val imageFilePath: String,
    val contentType: String = "image/jpeg",
    val localScanId: Long
)

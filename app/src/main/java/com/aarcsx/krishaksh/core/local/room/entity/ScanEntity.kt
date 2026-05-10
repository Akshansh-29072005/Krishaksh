package com.aarcsx.krishaksh.core.local.room.entity

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "scan_history")
data class ScanEntity(
    @PrimaryKey(autoGenerate = true)
    val id: Int = 0,
    val cropName: String,
    val imageUrl: String,
    val remoteScanId: String? = null,
    val predictionStatus: String = "PENDING",
    val diseaseName: String,
    val confidence: String,
    val symptoms: String,
    val prevention: String,
    val treatment: String,
    val recommendationTitle: String,
    val recommendationDesc: String,
    val capturedAt: Long
)

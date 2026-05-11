package com.aarcsx.krisho.features.scan

import androidx.compose.runtime.Immutable

@Immutable
data class ScanResult(
    val imageUrl: String,
    val diseaseName: String,
    val confidence: String,
    val symptoms: String,
    val prevention: String,
    val treatment: String,
    val recommendationTitle: String,
    val recommendationDesc: String
)

@Immutable
data class ScanUiState(
    val isAnalyzing: Boolean = false,
    val result: ScanResult? = null,
    val capturedImageUri: String? = null,
    val error: String? = null
)

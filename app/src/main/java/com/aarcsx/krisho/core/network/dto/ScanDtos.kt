package com.aarcsx.krisho.core.network.dto

data class PresignedUploadDto(val upload_url: String, val image_key: String)
data class CreateScanRequestDto(val crop_type: String, val image_key: String)
data class ScanDto(
    val id: String,
    val image_url: String,
    val crop_type: String,
    val prediction_status: String,
    val disease_name: String? = null,
    val confidence_score: Double? = null,
    val created_at: String
)
data class RecommendationDto(val disease: DiseaseDto? = null, val recommended_products: List<ProductDto> = emptyList())
data class DiseaseDto(val id: String, val name: String, val symptoms: String? = null, val prevention: String? = null, val treatment: String? = null)

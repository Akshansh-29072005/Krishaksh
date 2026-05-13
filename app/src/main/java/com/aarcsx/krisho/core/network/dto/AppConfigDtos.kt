package com.aarcsx.krisho.core.network.dto

data class CropDto(
    val name: String,
    val emoji: String
)

data class AppConfigDto(
    val minimum_version_code: Int,
    val latest_version_name: String? = null,
    val update_url: String? = null,
    val message: String? = null
)

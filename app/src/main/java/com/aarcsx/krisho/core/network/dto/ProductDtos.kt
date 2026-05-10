package com.aarcsx.krisho.core.network.dto

data class ProductDto(
    val id: String,
    val name: String,
    val description: String? = null,
    val crop_type: String? = null,
    val price: Double? = null,
    val currency: String? = "INR",
    val unit: String? = null,
    val image_url: String? = null
)

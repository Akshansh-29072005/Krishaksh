package com.aarcsx.krishaksh.core.local.room.entity

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "products")
data class ProductEntity(
    @PrimaryKey
    val id: String,
    val name: String,
    val price: Double,
    val imageUrl: String,
    val category: String,
    val description: String,
    val companyName: String,
    val usageInstructions: String,
    val lastUpdated: Long
)

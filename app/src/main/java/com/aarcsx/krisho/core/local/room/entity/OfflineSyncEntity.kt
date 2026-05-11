package com.aarcsx.krisho.core.local.room.entity

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "offline_sync_queue")
data class OfflineSyncEntity(
    @PrimaryKey(autoGenerate = true)
    val id: Int = 0,
    val type: String, // "SCAN_UPLOAD", "TICKET_CREATE", etc.
    val payload: String, // JSON payload string
    val createdAt: Long = System.currentTimeMillis()
)

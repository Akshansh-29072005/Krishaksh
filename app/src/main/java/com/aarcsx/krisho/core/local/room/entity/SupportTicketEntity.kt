package com.aarcsx.krisho.core.local.room.entity

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "support_tickets")
data class SupportTicketEntity(
    @PrimaryKey
    val ticketId: String,
    val subject: String,
    val status: String,
    val createdAt: Long,
    val lastUpdatedAt: Long
)

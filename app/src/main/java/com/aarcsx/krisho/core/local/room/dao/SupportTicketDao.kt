package com.aarcsx.krisho.core.local.room.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import com.aarcsx.krisho.core.local.room.entity.SupportTicketEntity
import kotlinx.coroutines.flow.Flow

@Dao
interface SupportTicketDao {
    @Query("SELECT * FROM support_tickets ORDER BY lastUpdatedAt DESC")
    fun getAllTickets(): Flow<List<SupportTicketEntity>>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insertTicket(ticket: SupportTicketEntity)
}

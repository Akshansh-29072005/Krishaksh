package com.aarcsx.krisho.core.local.room.dao

import androidx.room.*
import com.aarcsx.krisho.core.local.room.entity.OfflineSyncEntity
import kotlinx.coroutines.flow.Flow

@Dao
interface OfflineSyncDao {
    @Query("SELECT * FROM offline_sync_queue ORDER BY createdAt ASC")
    fun getAllPending(): Flow<List<OfflineSyncEntity>>

    @Query("SELECT * FROM offline_sync_queue ORDER BY createdAt ASC")
    suspend fun getAllPendingOnce(): List<OfflineSyncEntity>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun enqueue(entity: OfflineSyncEntity)

    @Delete
    suspend fun delete(entity: OfflineSyncEntity)

    @Query("DELETE FROM offline_sync_queue WHERE id = :id")
    suspend fun deleteById(id: Int)
}

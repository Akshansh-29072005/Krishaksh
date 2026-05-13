package com.aarcsx.krisho.core.local.room.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import androidx.room.Update
import com.aarcsx.krisho.core.local.room.entity.ScanEntity
import kotlinx.coroutines.flow.Flow

@Dao
interface ScanDao {
    @Query("SELECT * FROM scan_history ORDER BY capturedAt DESC")
    fun getAllScans(): Flow<List<ScanEntity>>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insertScan(scan: ScanEntity): Long

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insertScans(scans: List<ScanEntity>)

    @Query("SELECT * FROM scan_history WHERE remoteScanId = :remoteId LIMIT 1")
    suspend fun getByRemoteId(remoteId: String): ScanEntity?

    @Query("SELECT * FROM scan_history WHERE remoteScanId IS NULL OR predictionStatus = 'QUEUED' ORDER BY capturedAt ASC")
    suspend fun getPendingLocalScans(): List<ScanEntity>

    @Query("SELECT * FROM scan_history WHERE id = :id LIMIT 1")
    suspend fun getById(id: Int): ScanEntity?

    @Update
    suspend fun updateScan(scan: ScanEntity)

    @Query("DELETE FROM scan_history")
    suspend fun clearHistory()
}

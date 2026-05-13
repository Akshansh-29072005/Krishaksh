package com.aarcsx.krisho.core.repository

import com.aarcsx.krisho.core.local.room.dao.OfflineSyncDao
import com.aarcsx.krisho.core.local.room.entity.OfflineSyncEntity
import kotlinx.coroutines.flow.Flow
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class OfflineSyncRepository @Inject constructor(
    private val offlineSyncDao: OfflineSyncDao
) {
    fun getPendingTasks(): Flow<List<OfflineSyncEntity>> = offlineSyncDao.getAllPending()

    suspend fun getPendingTasksOnce(): List<OfflineSyncEntity> = offlineSyncDao.getAllPendingOnce()

    suspend fun enqueue(type: String, payload: String) {
        offlineSyncDao.enqueue(OfflineSyncEntity(type = type, payload = payload))
    }

    suspend fun resolve(id: Int) {
        offlineSyncDao.deleteById(id)
    }
}

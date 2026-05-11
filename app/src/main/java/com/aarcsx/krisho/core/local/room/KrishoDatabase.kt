package com.aarcsx.krisho.core.local.room

import androidx.room.Database
import androidx.room.RoomDatabase
import com.aarcsx.krisho.core.local.room.dao.OfflineSyncDao
import com.aarcsx.krisho.core.local.room.dao.ProductDao
import com.aarcsx.krisho.core.local.room.dao.ScanDao
import com.aarcsx.krisho.core.local.room.dao.SupportTicketDao
import com.aarcsx.krisho.core.local.room.entity.OfflineSyncEntity
import com.aarcsx.krisho.core.local.room.entity.ProductEntity
import com.aarcsx.krisho.core.local.room.entity.ScanEntity
import com.aarcsx.krisho.core.local.room.entity.SupportTicketEntity

@Database(
    entities = [
        ScanEntity::class,
        ProductEntity::class,
        SupportTicketEntity::class,
        OfflineSyncEntity::class
    ],
    version = 3,
    exportSchema = false
)
abstract class KrishoDatabase : RoomDatabase() {
    abstract fun scanDao(): ScanDao
    abstract fun productDao(): ProductDao
    abstract fun supportTicketDao(): SupportTicketDao
    abstract fun offlineSyncDao(): OfflineSyncDao
}

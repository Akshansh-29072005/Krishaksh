package com.aarcsx.krishaksh.core.local.room

import androidx.room.Database
import androidx.room.RoomDatabase
import com.aarcsx.krishaksh.core.local.room.dao.ProductDao
import com.aarcsx.krishaksh.core.local.room.dao.ScanDao
import com.aarcsx.krishaksh.core.local.room.dao.SupportTicketDao
import com.aarcsx.krishaksh.core.local.room.entity.ProductEntity
import com.aarcsx.krishaksh.core.local.room.entity.ScanEntity
import com.aarcsx.krishaksh.core.local.room.entity.SupportTicketEntity

@Database(
    entities = [
        ScanEntity::class,
        ProductEntity::class,
        SupportTicketEntity::class
    ],
    version = 2,
    exportSchema = false
)
abstract class KrishakshDatabase : RoomDatabase() {
    abstract fun scanDao(): ScanDao
    abstract fun productDao(): ProductDao
    abstract fun supportTicketDao(): SupportTicketDao
}

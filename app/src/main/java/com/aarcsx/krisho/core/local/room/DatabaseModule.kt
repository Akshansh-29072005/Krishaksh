package com.aarcsx.krisho.core.local.room

import android.content.Context
import androidx.room.Room
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
object DatabaseModule {

    @Provides
    @Singleton
    fun provideKrishoDatabase(@ApplicationContext context: Context): KrishoDatabase {
        return Room.databaseBuilder(
            context,
            KrishoDatabase::class.java,
            "krisho_db"
        )
        .fallbackToDestructiveMigration()
        .build()
    }

    @Provides
    fun provideScanDao(db: KrishoDatabase) = db.scanDao()

    @Provides
    fun provideProductDao(db: KrishoDatabase) = db.productDao()

    @Provides
    fun provideSupportTicketDao(db: KrishoDatabase) = db.supportTicketDao()
}

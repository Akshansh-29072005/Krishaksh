package com.aarcsx.krishaksh.core.local.room

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
    fun provideKrishakshDatabase(@ApplicationContext context: Context): KrishakshDatabase {
        return Room.databaseBuilder(
            context,
            KrishakshDatabase::class.java,
            "krishaksh_db"
        )
        .fallbackToDestructiveMigration()
        .build()
    }

    @Provides
    fun provideScanDao(db: KrishakshDatabase) = db.scanDao()

    @Provides
    fun provideProductDao(db: KrishakshDatabase) = db.productDao()

    @Provides
    fun provideSupportTicketDao(db: KrishakshDatabase) = db.supportTicketDao()
}

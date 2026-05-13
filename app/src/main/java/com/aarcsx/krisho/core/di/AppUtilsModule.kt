package com.aarcsx.krisho.core.di

import android.content.Context
import com.aarcsx.krisho.core.auth.GoogleSignInManager
import com.aarcsx.krisho.core.util.LocationProvider
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
object AppUtilsModule {

    @Provides
    @Singleton
    fun provideLocationProvider(@ApplicationContext context: Context): LocationProvider {
        return LocationProvider(context)
    }

    @Provides
    @Singleton
    fun provideGoogleSignInManager(@ApplicationContext context: Context): GoogleSignInManager {
        return GoogleSignInManager(context)
    }

    @Provides
    @Singleton
    fun provideNetworkMonitor(@ApplicationContext context: Context): com.aarcsx.krisho.core.util.NetworkMonitor {
        return com.aarcsx.krisho.core.util.NetworkMonitor(context)
    }
}

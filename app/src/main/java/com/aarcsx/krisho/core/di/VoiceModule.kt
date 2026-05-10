package com.aarcsx.krisho.core.di

import android.content.Context
import com.aarcsx.krisho.core.common.voice.VoiceRecorderManager
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
object VoiceModule {

    @Provides
    @Singleton
    fun provideVoiceRecorderManager(
        @ApplicationContext context: Context
    ): VoiceRecorderManager {
        return VoiceRecorderManager(context)
    }
}
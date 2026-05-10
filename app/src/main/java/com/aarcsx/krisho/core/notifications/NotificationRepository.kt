package com.aarcsx.krisho.core.notifications

import android.util.Log
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class NotificationRepository @Inject constructor() {
    suspend fun fetchFcmToken(): String? = null

    suspend fun syncFcmTokenWithBackend(token: String) {
        Log.d("NotificationRepository", "FCM token ready for backend sync: $token")
    }
}

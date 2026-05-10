package com.aarcsx.krisho.core.utils

import android.app.ActivityManager
import android.content.Context
import android.net.ConnectivityManager
import android.net.NetworkCapabilities
import timber.log.Timber

/**
 * Utilities for real-device validation and resilience testing.
 */
object SystemHealthMonitor {

    /**
     * Checks if the device has low memory.
     */
    fun checkMemoryStatus(context: Context) {
        val activityManager = context.getSystemService(Context.ACTIVITY_SERVICE) as ActivityManager
        val memoryInfo = ActivityManager.MemoryInfo()
        activityManager.getMemoryInfo(memoryInfo)
        
        if (memoryInfo.lowMemory) {
            Timber.w("Device is in LOW MEMORY state.")
        }
        Timber.d("Available Memory: ${memoryInfo.availMem / (1024 * 1024)} MB")
    }

    /**
     * Checks network type for poor connection validation.
     */
    fun getNetworkQuality(context: Context): String {
        val connectivityManager = context.getSystemService(Context.CONNECTIVITY_SERVICE) as ConnectivityManager
        val network = connectivityManager.activeNetwork ?: return "None"
        val caps = connectivityManager.getNetworkCapabilities(network) ?: return "None"
        
        return when {
            caps.hasTransport(NetworkCapabilities.TRANSPORT_WIFI) -> "WiFi (Fast)"
            caps.hasTransport(NetworkCapabilities.TRANSPORT_CELLULAR) -> {
                // Simplified assessment
                "Cellular"
            }
            else -> "Other"
        }
    }
}
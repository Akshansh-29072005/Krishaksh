package com.aarcsx.krisho.core.util

import android.Manifest
import android.annotation.SuppressLint
import android.content.Context
import android.content.pm.PackageManager
import android.location.Location
import android.location.LocationManager
import android.util.Log
import androidx.core.app.ActivityCompat
import com.google.android.gms.location.FusedLocationProviderClient
import com.google.android.gms.location.LocationServices
import com.google.android.gms.location.Priority
import com.google.android.gms.tasks.CancellationTokenSource
import dagger.hilt.android.qualifiers.ApplicationContext
import kotlinx.coroutines.suspendCancellableCoroutine
import javax.inject.Inject
import javax.inject.Singleton
import kotlin.coroutines.resume

private const val TAG = "LocationProvider"

@Singleton
class LocationProvider @Inject constructor(
    @ApplicationContext private val context: Context
) {
    private val fusedLocationClient: FusedLocationProviderClient =
        LocationServices.getFusedLocationProviderClient(context)

    @SuppressLint("MissingPermission")
    suspend fun getCurrentLocation(): Location? {
        if (!hasLocationPermission()) {
            Log.w(TAG, "Location permission NOT granted")
            return null
        }

        if (!isDeviceLocationEnabled()) {
            Log.w(TAG, "Device location services are OFF")
            return null
        }

        Log.d(TAG, "Permission OK, device location ON. Fetching location...")

        return suspendCancellableCoroutine { continuation ->
            val cts = CancellationTokenSource()
            var resumed = false

            fun safeResume(location: Location?) {
                if (!resumed) {
                    resumed = true
                    Log.d(TAG, "Location result: lat=${location?.latitude}, lon=${location?.longitude}")
                    continuation.resume(location)
                }
            }

            // Try last known location first (instant)
            fusedLocationClient.lastLocation.addOnSuccessListener { lastLocation ->
                if (lastLocation != null) {
                    Log.d(TAG, "Got lastLocation: ${lastLocation.latitude}, ${lastLocation.longitude}")
                    safeResume(lastLocation)
                } else {
                    Log.d(TAG, "No lastLocation, requesting current...")
                    fusedLocationClient.getCurrentLocation(
                        Priority.PRIORITY_BALANCED_POWER_ACCURACY,
                        cts.token
                    ).addOnSuccessListener { location: Location? ->
                        Log.d(TAG, "getCurrentLocation result: $location")
                        safeResume(location)
                    }.addOnFailureListener { e ->
                        Log.e(TAG, "getCurrentLocation failed: ${e.message}")
                        safeResume(null)
                    }
                }
            }.addOnFailureListener { e ->
                Log.e(TAG, "lastLocation failed: ${e.message}")
                safeResume(null)
            }

            continuation.invokeOnCancellation {
                cts.cancel()
            }
        }
    }

    private fun hasLocationPermission(): Boolean {
        return ActivityCompat.checkSelfPermission(
            context,
            Manifest.permission.ACCESS_COARSE_LOCATION
        ) == PackageManager.PERMISSION_GRANTED || ActivityCompat.checkSelfPermission(
            context,
            Manifest.permission.ACCESS_FINE_LOCATION
        ) == PackageManager.PERMISSION_GRANTED
    }

    private fun isDeviceLocationEnabled(): Boolean {
        val lm = context.getSystemService(Context.LOCATION_SERVICE) as LocationManager
        return lm.isProviderEnabled(LocationManager.GPS_PROVIDER) ||
                lm.isProviderEnabled(LocationManager.NETWORK_PROVIDER)
    }
}

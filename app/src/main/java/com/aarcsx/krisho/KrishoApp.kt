package com.aarcsx.krisho

import android.app.Application
import android.util.Log
import com.aarcsx.krisho.core.local.datastore.PreferencesManager
import dagger.hilt.android.HiltAndroidApp
import timber.log.Timber
import javax.inject.Inject

@HiltAndroidApp
class KrishoApp : Application() {
    @Inject
    lateinit var preferencesManager: PreferencesManager

    override fun onCreate() {
        super.onCreate()
        
        if (BuildConfig.DEBUG) {
            Timber.plant(Timber.DebugTree())
        } else {
            Timber.plant(ReleaseTree())
        }
    }

    /**
     * Custom Timber Tree for release builds to send important logs to Crashlytics
     */
    private class ReleaseTree : Timber.Tree() {
        override fun log(priority: Int, tag: String?, message: String, t: Throwable?) {
            if (priority == Log.VERBOSE || priority == Log.DEBUG) {
                return
            }
            // Add Firebase Crashlytics logging here
            // FirebaseCrashlytics.getInstance().log(message)
            // t?.let { FirebaseCrashlytics.getInstance().recordException(it) }
        }
    }
}

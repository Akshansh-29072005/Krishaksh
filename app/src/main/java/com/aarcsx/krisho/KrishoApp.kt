package com.aarcsx.krisho

import android.app.Application
import android.util.Log
import com.aarcsx.krisho.core.local.datastore.PreferencesManager
import com.google.firebase.crashlytics.FirebaseCrashlytics
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
            val crashlytics = FirebaseCrashlytics.getInstance()
            crashlytics.log(if (tag != null) "[$tag] $message" else message)
            if (t != null) {
                crashlytics.recordException(t)
            }
        }
    }
}

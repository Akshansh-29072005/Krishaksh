package com.aarcsx.krishaksh.core.common

import android.content.Context
import android.content.ContextWrapper
import android.content.res.Configuration
import android.os.LocaleList
import com.aarcsx.krishaksh.core.local.datastore.PreferencesManager
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.runBlocking
import java.util.Locale

object LocaleManager {
    fun updateResources(context: Context, language: String): Context {
        val locale = Locale(language)
        Locale.setDefault(locale)

        val res = context.resources
        val config = Configuration(res.configuration)
        
        config.setLocale(locale)
        val localeList = LocaleList(locale)
        config.setLocales(localeList)
        
        return context.createConfigurationContext(config)
    }

    fun getLocaleContextWrapper(context: Context, preferencesManager: PreferencesManager): ContextWrapper {
        val language = runBlocking { preferencesManager.languageSetting.first() }
        return ContextWrapper(updateResources(context, language))
    }
}

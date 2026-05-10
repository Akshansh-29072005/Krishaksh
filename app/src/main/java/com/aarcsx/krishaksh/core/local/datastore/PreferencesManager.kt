package com.aarcsx.krishaksh.core.local.datastore

import android.content.Context
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.*
import androidx.datastore.preferences.preferencesDataStore
import dagger.hilt.android.qualifiers.ApplicationContext
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.catch
import kotlinx.coroutines.flow.map
import java.io.IOException
import javax.inject.Inject
import javax.inject.Singleton

val Context.dataStore: DataStore<Preferences> by preferencesDataStore(name = "krishaksh_prefs")

@Singleton
class PreferencesManager @Inject constructor(
    @ApplicationContext private val context: Context
) {
    private val dataStore = context.dataStore

    companion object {
        val JWT_TOKEN = stringPreferencesKey("jwt_token")
        val REFRESH_TOKEN = stringPreferencesKey("refresh_token")
        val LANGUAGE_SETTING = stringPreferencesKey("language_setting")
        val IS_FIRST_LAUNCH = booleanPreferencesKey("is_first_launch")
    }

    val jwtToken: Flow<String?> = dataStore.data
        .catch { exception ->
            if (exception is IOException) emit(emptyPreferences()) else throw exception
        }
        .map { it[JWT_TOKEN] }

    val refreshToken: Flow<String?> = dataStore.data
        .catch { exception ->
            if (exception is IOException) emit(emptyPreferences()) else throw exception
        }
        .map { it[REFRESH_TOKEN] }

    val languageSetting: Flow<String> = dataStore.data
        .catch { exception ->
            if (exception is IOException) emit(emptyPreferences()) else throw exception
        }
        .map { it[LANGUAGE_SETTING] ?: "en" }

    suspend fun saveTokens(jwt: String, refresh: String) {
        dataStore.edit { prefs ->
            prefs[JWT_TOKEN] = jwt
            prefs[REFRESH_TOKEN] = refresh
        }
    }

    suspend fun clearTokens() {
        dataStore.edit { prefs ->
            prefs.remove(JWT_TOKEN)
            prefs.remove(REFRESH_TOKEN)
        }
    }

    suspend fun saveLanguage(lang: String) {
        dataStore.edit { prefs ->
            prefs[LANGUAGE_SETTING] = lang
        }
    }
}

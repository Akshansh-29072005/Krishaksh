package com.aarcsx.krisho.core.local.datastore

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
import com.aarcsx.krisho.core.network.dto.AppConfigDto

val Context.dataStore: DataStore<Preferences> by preferencesDataStore(name = "krisho_prefs")

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
        val FCM_TOKEN = stringPreferencesKey("fcm_token")
        val MINIMUM_VERSION_CODE = intPreferencesKey("minimum_version_code")
        val LATEST_VERSION_NAME = stringPreferencesKey("latest_version_name")
        val APP_UPDATE_URL = stringPreferencesKey("app_update_url")
        val APP_CONFIG_MESSAGE = stringPreferencesKey("app_config_message")
    }

    val fcmToken: Flow<String?> = dataStore.data
        .catch { exception ->
            if (exception is IOException) emit(emptyPreferences()) else throw exception
        }
        .map { it[FCM_TOKEN] }

    suspend fun saveFcmToken(token: String) {
        dataStore.edit { prefs ->
            prefs[FCM_TOKEN] = token
        }
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

    val cachedAppConfig: Flow<AppConfigDto?> = dataStore.data
        .catch { exception ->
            if (exception is IOException) emit(emptyPreferences()) else throw exception
        }
        .map { prefs ->
            prefs[MINIMUM_VERSION_CODE]?.let { minimumVersion ->
                AppConfigDto(
                    minimum_version_code = minimumVersion,
                    latest_version_name = prefs[LATEST_VERSION_NAME],
                    update_url = prefs[APP_UPDATE_URL],
                    message = prefs[APP_CONFIG_MESSAGE]
                )
            }
        }

    suspend fun saveAppConfig(config: AppConfigDto) {
        dataStore.edit { prefs ->
            prefs[MINIMUM_VERSION_CODE] = config.minimum_version_code
            config.latest_version_name?.let { prefs[LATEST_VERSION_NAME] = it } ?: prefs.remove(LATEST_VERSION_NAME)
            config.update_url?.let { prefs[APP_UPDATE_URL] = it } ?: prefs.remove(APP_UPDATE_URL)
            config.message?.let { prefs[APP_CONFIG_MESSAGE] = it } ?: prefs.remove(APP_CONFIG_MESSAGE)
        }
    }

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

package com.aarcsx.krisho.core.auth

import com.aarcsx.krisho.core.local.datastore.PreferencesManager
import kotlinx.coroutines.flow.first
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class TokenManager @Inject constructor(
    private val prefs: PreferencesManager
) {
    suspend fun accessToken(): String? = prefs.jwtToken.first()
    suspend fun refreshToken(): String? = prefs.refreshToken.first()
    suspend fun save(access: String, refresh: String) = prefs.saveTokens(access, refresh)
    suspend fun clear() = prefs.clearTokens()
}

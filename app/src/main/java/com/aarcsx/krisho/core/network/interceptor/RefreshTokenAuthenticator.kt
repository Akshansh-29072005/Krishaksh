package com.aarcsx.krisho.core.network.interceptor

import com.aarcsx.krisho.core.auth.TokenManager
import com.aarcsx.krisho.core.network.api.AuthApiService
import com.aarcsx.krisho.core.network.dto.RefreshTokenRequest
import kotlinx.coroutines.runBlocking
import okhttp3.Authenticator
import okhttp3.Request
import okhttp3.Response
import okhttp3.Route
import javax.inject.Inject

class RefreshTokenAuthenticator @Inject constructor(
    private val tokenManager: TokenManager,
    private val authApi: dagger.Lazy<AuthApiService>
) : Authenticator {
    override fun authenticate(route: Route?, response: Response): Request? {
        if (responseCount(response) >= 2) return null
        val refresh = runBlocking { tokenManager.refreshToken() } ?: return null
        val refreshed = runBlocking {
            val res = authApi.get().refresh(RefreshTokenRequest(refresh))
            val body = res.body()?.data
            if (res.isSuccessful && body != null) {
                tokenManager.save(body.access_token, body.refresh_token)
                body.access_token
            } else null
        } ?: return null
        return response.request.newBuilder()
            .header("Authorization", "Bearer $refreshed")
            .build()
    }

    private fun responseCount(response: Response): Int {
        var count = 1
        var prior = response.priorResponse
        while (prior != null) {
            count++
            prior = prior.priorResponse
        }
        return count
    }
}

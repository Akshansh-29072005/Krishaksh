package com.aarcsx.krishaksh.core.network.interceptor

import com.aarcsx.krishaksh.core.auth.TokenManager
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.firstOrNull
import kotlinx.coroutines.runBlocking
import kotlinx.coroutines.withTimeoutOrNull
import okhttp3.Interceptor
import okhttp3.Response
import javax.inject.Inject

class AuthInterceptor @Inject constructor(
    private val tokenManager: TokenManager
) : Interceptor {
    override fun intercept(chain: Interceptor.Chain): Response {
        val request = chain.request()
        
        // Skip auth for login/refresh/ads
        if (request.url.encodedPath.contains("auth/") || request.url.encodedPath.contains("ads")) {
            return chain.proceed(request)
        }

        val token = runBlocking { 
            withTimeoutOrNull(2000) { tokenManager.accessToken() }
        }
        val req = request.newBuilder().apply {
            if (!token.isNullOrBlank()) addHeader("Authorization", "Bearer $token")
            addHeader("Accept", "application/json")
        }.build()
        return try {
            chain.proceed(req)
        } catch (e: Exception) {
            throw e
        }
    }
}

package com.aarcsx.krisho.core.network.interceptor

import com.aarcsx.krisho.core.auth.TokenManager
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

        // Skip auth for login/refresh/ads endpoints
        if (request.url.encodedPath.contains("auth/") || request.url.encodedPath.contains("ads")) {
            return chain.proceed(request)
        }

        // Skip adding Authorization header for AWS S3 presigned URL uploads (auth is in query parameters)
        if (request.url.host.contains("amazonaws.com")) {
            // Still add Accept header for consistency
            val noAuthReq = request.newBuilder()
                .addHeader("Accept", "application/json")
                .build()
            return chain.proceed(noAuthReq)
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

package com.aarcsx.krisho.core.network.dto

data class GoogleLoginRequest(val id_token: String)
data class RefreshTokenRequest(val refresh_token: String)
data class AuthTokensDto(val access_token: String, val refresh_token: String, val expires_in: Int)
data class LogoutRequest(val refresh_token: String)

package com.aarcsx.krishaksh.core.network

/**
 * Defines all the backend API endpoints and standard network configuration.
 */
object ApiEndpoints {

    // Base URL for the Go Gin Backend
    // Using 10.0.2.2 because it is the standard loopback IP for Android Emulator to hit localhost.
    const val BASE_URL = "http://10.0.2.2:8080/api/v1/"
}

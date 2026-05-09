package com.aarcsx.krishaksh.core.network

/**
 * Defines all the backend API endpoints and standard network configuration.
 */
object ApiEndpoints {

    // Base URL for the Go Gin Backend
    // Using 10.0.2.2 because it is the standard loopback IP for Android Emulator to hit localhost.
    const val BASE_URL = "http://10.0.2.2:8080/api/v1/"

    // Weather & Alerts
    const val GET_CURRENT_WEATHER = "weather/current"
    const val GET_WEATHER_ALERTS = "weather/alerts"
    
    // AI Scanning
    // Expected to accept Multipart Form Data: image file & crop name
    const val POST_SCAN_CROP = "scan/analyze"
    
    // User Scan History
    const val GET_RECENT_SCANS = "user/scans"
    
    // Marketplace Products
    const val GET_ALL_PRODUCTS = "products"
    // Append the product ID dynamically, e.g., "products/123"
    const val GET_PRODUCT_DETAILS = "products/{id}"
    
    // Auth & Profile (Placeholders for future)
    const val POST_LOGIN = "auth/login"
    const val POST_REGISTER = "auth/register"
    const val GET_USER_PROFILE = "user/profile"
}
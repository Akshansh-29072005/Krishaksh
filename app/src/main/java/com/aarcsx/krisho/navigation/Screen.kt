package com.aarcsx.krisho.navigation

sealed class Screen(val route: String) {
    object Auth : Screen("auth")
    object Home : Screen("home")
    object History : Screen("history")
    object Products : Screen("products")
    object Support : Screen("support")
    object Profile : Screen("profile")
    object Scan : Screen("scan")
    object Result : Screen("result")
    object Loading : Screen("loading")
    object ProductDetail : Screen("product_detail/{productId}") {
        fun createRoute(productId: String) = "product_detail/$productId"
    }
    object HelpCenter : Screen("help_center")
    object PrivacyPolicy : Screen("privacy_policy")
    object TermsOfService : Screen("terms_of_service")
}

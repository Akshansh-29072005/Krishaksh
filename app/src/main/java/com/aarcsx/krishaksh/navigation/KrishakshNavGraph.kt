package com.aarcsx.krishaksh.navigation

import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.navigation.NavHostController
import androidx.navigation.compose.*
import androidx.navigation.NavType
import androidx.navigation.navArgument
import androidx.lifecycle.viewmodel.compose.viewModel

import com.aarcsx.krishaksh.features.home.HomeScreen
import com.aarcsx.krishaksh.features.history.HistoryScreen
import com.aarcsx.krishaksh.features.products.ProductsScreen
import com.aarcsx.krishaksh.features.products.ProductDetailScreen
import com.aarcsx.krishaksh.features.support.SupportScreen
import com.aarcsx.krishaksh.features.profile.ProfileScreen
import com.aarcsx.krishaksh.features.scan.*

data class BottomNavItem(
    val screen: Screen,
    val icon: ImageVector,
    val label: String
)

val bottomNavItems = listOf(
    BottomNavItem(Screen.Home, Icons.Default.Home, "Home"),
    BottomNavItem(Screen.History, Icons.Default.History, "History"),
    BottomNavItem(Screen.Products, Icons.Default.ShoppingCart, "Products"),
    BottomNavItem(Screen.Support, Icons.Default.SupportAgent, "Support"),
    BottomNavItem(Screen.Profile, Icons.Default.Person, "Profile")
)


@Composable
fun KrishakshNavGraph(
    navController: NavHostController = rememberNavController()
) {
    val navBackStackEntry by navController.currentBackStackEntryAsState()
    val currentRoute = navBackStackEntry?.destination?.route

    Scaffold(
        bottomBar = {
            if (currentRoute in listOf("home", "history", "products", "support", "profile")) {
                NavigationBar(
                    containerColor = Color.White,
                    contentColor = Color(0xFF1F4D36)
                ) {
                    bottomNavItems.forEach { item ->
                        NavigationBarItem(
                            icon = { Icon(item.icon, contentDescription = item.label) },
                            label = { Text(item.label) },
                            selected = currentRoute == item.screen.route,
                            onClick = {
                                navController.navigate(item.screen.route) {
                                    popUpTo(navController.graph.startDestinationId) {
                                        saveState = true
                                    }
                                    launchSingleTop = true
                                    restoreState = true
                                }
                            },
                            colors = NavigationBarItemDefaults.colors(
                                selectedIconColor = Color(0xFF1F4D36),
                                unselectedIconColor = Color.Gray,
                                selectedTextColor = Color(0xFF1F4D36),
                                unselectedTextColor = Color.Gray,
                                indicatorColor = Color(0xFFE8F5E9)
                            )
                        )
                    }
                }
            }
        }
    ) { innerPadding ->
        NavHost(
            navController = navController,
            startDestination = Screen.Home.route,
            modifier = Modifier.padding(innerPadding)
        ) {
            composable(Screen.Home.route) {
                HomeScreen(onScanClick = { navController.navigate(Screen.Scan.route) })
            }
            composable(Screen.History.route) {
                HistoryScreen(onBackClick = { navController.popBackStack() })
            }
            composable(Screen.Products.route) {
                ProductsScreen(
                    onBackClick = { navController.popBackStack() },
                    onProductClick = { productId ->
                        navController.navigate(Screen.ProductDetail.createRoute(productId))
                    }
                )
            }
            composable(
                route = Screen.ProductDetail.route,
                arguments = listOf(navArgument("productId") { type = NavType.StringType })
            ) { backStackEntry ->
                val productId = backStackEntry.arguments?.getString("productId") ?: ""
                ProductDetailScreen(
                    productId = productId,
                    onBackClick = { navController.popBackStack() }
                )
            }
            composable(Screen.Support.route) {
                SupportScreen(onBackClick = { navController.popBackStack() })
            }
            composable(Screen.Profile.route) {
                ProfileScreen(onBackClick = { navController.popBackStack() })
            }
            composable(Screen.Scan.route) {
                ScanScreen(
                    onBackClick = { navController.popBackStack() },
                    onCaptured = { 
                        navController.navigate(Screen.Loading.route)
                    }
                )
            }
            composable(Screen.Loading.route) {
                LaunchedEffect(Unit) {
                    kotlinx.coroutines.delay(2000)
                    navController.navigate(Screen.Result.route) {
                        popUpTo(Screen.Loading.route) { inclusive = true }
                    }
                }
                LoadingScreen()
            }
            composable(Screen.Result.route) {
                val mockResult = ScanResult(
                    diseaseName = "Potato Late Blight",
                    confidence = "98%",
                    symptoms = "Dark, water-soaked spots on leaves that enlarge rapidly.",
                    prevention = "Use certified disease-free seeds and rotate crops.",
                    treatment = "Apply fungicide immediately. Remove infected plants.",
                    recommendationTitle = "Copper Fungicide",
                    recommendationDesc = "Effective against late blight in potatoes and tomatoes."
                )
                ResultScreen(
                    result = mockResult,
                    onBackClick = { navController.popBackStack(Screen.Home.route, false) }
                )
            }
        }
    }
}

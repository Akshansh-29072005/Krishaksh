package com.aarcsx.krisho.navigation

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
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.ui.Alignment
import com.aarcsx.krisho.core.designsystem.components.ForestGreen

import com.aarcsx.krisho.features.home.HomeScreen
import com.aarcsx.krisho.features.auth.AuthScreen
import com.aarcsx.krisho.features.history.HistoryScreen
import com.aarcsx.krisho.features.products.ProductsScreen
import com.aarcsx.krisho.features.products.ProductDetailScreen
import com.aarcsx.krisho.features.support.SupportScreen
import com.aarcsx.krisho.features.profile.ProfileScreen
import com.aarcsx.krisho.features.scan.*
import com.aarcsx.krisho.features.legal.*

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
fun KrishoNavGraph(
    navController: NavHostController = rememberNavController(),
    startDestination: String = Screen.Auth.route
) {
    val navBackStackEntry by navController.currentBackStackEntryAsState()
    val currentRoute = navBackStackEntry?.destination?.route

    NavHost(
        navController = navController,
        startDestination = startDestination
    ) {
        composable(Screen.Auth.route) {
            AuthScreen(onAuthSuccess = {
                navController.navigate(Screen.Home.route) {
                    popUpTo(Screen.Auth.route) { inclusive = true }
                }
            })
        }
        composable(Screen.Home.route) {
            MainScaffold(navController) {
                HomeScreen(onScanClick = { navController.navigate(Screen.CropSelection.route) })
            }
        }
        composable(Screen.History.route) {
            MainScaffold(navController) {
                HistoryScreen(onBackClick = { navController.popBackStack() })
            }
        }
        composable(Screen.Products.route) {
            MainScaffold(navController) {
                ProductsScreen(
                    onBackClick = { navController.popBackStack() },
                    onProductClick = { productId ->
                        navController.navigate(Screen.ProductDetail.createRoute(productId))
                    }
                )
            }
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
            MainScaffold(navController) {
                SupportScreen(onBackClick = { navController.popBackStack() })
            }
        }
        composable(Screen.Profile.route) {
            MainScaffold(navController) {
                ProfileScreen(
                    onBackClick = { navController.popBackStack() },
                    onNavigateToHelp = { navController.navigate(Screen.HelpCenter.route) },
                    onNavigateToPrivacy = { navController.navigate(Screen.PrivacyPolicy.route) },
                    onNavigateToTerms = { navController.navigate(Screen.TermsOfService.route) },
                    onLogoutSuccess = {
                        navController.navigate(Screen.Auth.route) {
                            popUpTo(Screen.Home.route) { inclusive = true }
                        }
                    }
                )
            }
        }
        composable(Screen.HelpCenter.route) {
            HelpCenterScreen(onBackClick = { navController.popBackStack() })
        }
        composable(Screen.PrivacyPolicy.route) {
            PrivacyPolicyScreen(onBackClick = { navController.popBackStack() })
        }
        composable(Screen.TermsOfService.route) {
            TermsOfServiceScreen(onBackClick = { navController.popBackStack() })
        }
        composable(Screen.CropSelection.route) {
            CropSelectionScreen(
                onBackClick = { navController.popBackStack() },
                onCropSelected = { cropName ->
                    navController.navigate(Screen.Scan.createRoute(cropName))
                }
            )
        }
        composable(
            route = Screen.Scan.route,
            arguments = listOf(navArgument("cropName") { type = NavType.StringType })
        ) { backStackEntry ->
            val cropName = backStackEntry.arguments?.getString("cropName") ?: ""
            ScanScreen(
                cropName = cropName,
                onBackClick = { navController.popBackStack() },
                onCaptured = { result ->
                    navController.navigate(Screen.Result.route)
                }
            )
        }
        composable(Screen.Result.route) {
            val scanViewModel: ScanViewModel = hiltViewModel()
            val uiState by scanViewModel.uiState.collectAsState()
            
            if (uiState.result != null) {
                ResultScreen(
                    result = uiState.result!!,
                    onBackClick = { navController.popBackStack(Screen.Home.route, false) }
                )
            } else {
                Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
                    Column(horizontalAlignment = Alignment.CenterHorizontally) {
                        CircularProgressIndicator(color = ForestGreen)
                        if (uiState.error != null) {
                            Text(uiState.error!!, color = Color.Red)
                            Button(onClick = { navController.popBackStack() }) {
                                Text("Go Back")
                            }
                        }
                    }
                }
            }
        }
    }
}

@Composable
fun MainScaffold(
    navController: NavHostController,
    content: @Composable (androidx.compose.foundation.layout.PaddingValues) -> Unit
) {
    val navBackStackEntry by navController.currentBackStackEntryAsState()
    val currentRoute = navBackStackEntry?.destination?.route

    Scaffold(
        bottomBar = {
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
    ) { innerPadding ->
        content(innerPadding)
    }
}

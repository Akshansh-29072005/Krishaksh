import java.util.Properties

plugins {
    alias(libs.plugins.android.application)
    alias(libs.plugins.kotlin.android)
    alias(libs.plugins.compose.compiler)
    alias(libs.plugins.hilt)
    alias(libs.plugins.google.services)
    alias(libs.plugins.firebase.crashlytics)
    id("kotlin-kapt")
}

// Load signing credentials from local.properties (never committed to git)
val localProperties = Properties().apply {
    val localPropertiesFile = rootProject.file("local.properties")
    if (localPropertiesFile.exists()) {
        localPropertiesFile.inputStream().use { load(it) }
    }
}

android {
    namespace = "com.aarcsx.krisho"
    compileSdk = 35

    defaultConfig {
        applicationId = "com.aarcsx.krisho"
        minSdk = 24
        targetSdk = 35
        versionCode = 1
        versionName = "1.03"

        testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
    }

    signingConfigs {
        create("release") {
            storeFile = file(localProperties.getProperty("KEYSTORE_FILE", "../krisho-prod-key.jks"))
            storePassword = localProperties.getProperty("KEYSTORE_PASSWORD", "")
            keyAlias = localProperties.getProperty("KEY_ALIAS", "")
            keyPassword = localProperties.getProperty("KEY_PASSWORD", "")
        }
    }

    buildTypes {
        debug {
            signingConfig = signingConfigs.getByName("release")
            buildConfigField("String", "BASE_URL", "\"https://api-krisho.aarcsx.com/api/v1/\"")
            buildConfigField("String", "RAZORPAY_KEY_ID", "\"rzp_test_SngDkzeaoqQi4o\"")
            manifestPlaceholders["appLabel"] = "Krisho Dev"
        }
        create("staging") {
            initWith(getByName("release"))
            applicationIdSuffix = ".staging"
            matchingFallbacks.add("release")
            buildConfigField("String", "BASE_URL", "\"https://api.staging.krisho.aarcsx.com/api/v1/\"")
            buildConfigField("String", "RAZORPAY_KEY_ID", "\"rzp_test_SngDkzeaoqQi4o\"")
            manifestPlaceholders["appLabel"] = "Krisho Staging"
        }
        release {
            isMinifyEnabled = true
            isShrinkResources = true
            signingConfig = signingConfigs.getByName("release")
            buildConfigField("String", "BASE_URL", "\"https://api-krisho.aarcsx.com/api/v1/\"")
            buildConfigField("String", "RAZORPAY_KEY_ID", "\"${localProperties.getProperty("RAZORPAY_LIVE_KEY", "rzp_test_SngDkzeaoqQi4o")}\"")
            manifestPlaceholders["appLabel"] = "Krisho: Farmer's AI"
            proguardFiles(
                getDefaultProguardFile("proguard-android-optimize.txt"),
                "proguard-rules.pro"
            )
        }
    }
    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_11
        targetCompatibility = JavaVersion.VERSION_11
    }
    kotlinOptions {
        jvmTarget = "11"
    }
    buildFeatures {
        compose = true
        buildConfig = true
    }
}

dependencies {
    // Hilt
    implementation(libs.hilt.android)
    kapt(libs.hilt.compiler)
    implementation(libs.androidx.hilt.navigation.compose)

    implementation(libs.androidx.core.ktx)
    implementation(libs.androidx.lifecycle.runtime.ktx)
    implementation(libs.androidx.activity.compose)
    implementation(platform(libs.androidx.compose.bom))
    implementation(libs.androidx.ui)
    implementation(libs.androidx.ui.text)
    implementation(libs.androidx.ui.graphics)
    implementation(libs.androidx.ui.tooling.preview)
    implementation(libs.androidx.material3)
    implementation(libs.androidx.navigation.compose)
    implementation(libs.material)
    implementation(libs.androidx.material.icons.extended)

    // Firebase
    implementation(platform(libs.firebase.bom))
    implementation(libs.firebase.messaging)
    implementation(libs.firebase.analytics)
    implementation(libs.firebase.crashlytics)

    // Razorpay
    implementation(libs.razorpay)

    // Logging
    implementation(libs.timber)

    // Room
    implementation(libs.room.runtime)
    implementation(libs.room.ktx)
    kapt(libs.room.compiler)

    // Retrofit
    implementation(libs.retrofit)
    implementation(libs.retrofit.converter.gson)
    implementation(libs.okhttp.logging)
    implementation(libs.play.services.auth)
    implementation(libs.play.services.location)

    // Images
    implementation(libs.coil.compose)

    // CameraX
    implementation(libs.androidx.camera.core)
    implementation(libs.androidx.camera.camera2)
    implementation(libs.androidx.camera.lifecycle)
    implementation(libs.androidx.camera.view)

    // WorkManager
    implementation(libs.androidx.work.runtime.ktx)

    // DataStore
    implementation(libs.androidx.datastore.preferences)

    testImplementation(libs.junit)
    androidTestImplementation(libs.androidx.junit)
    androidTestImplementation(libs.androidx.espresso.core)
    androidTestImplementation(platform(libs.androidx.compose.bom))
    androidTestImplementation(libs.androidx.ui.test.junit4)
    debugImplementation(libs.androidx.ui.tooling)
    debugImplementation(libs.androidx.ui.test.manifest)
}

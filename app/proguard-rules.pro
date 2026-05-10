# Add project specific ProGuard rules here.
# You can control the set of applied configuration files using the
# proguardFiles setting in build.gradle.

# Hilt/Dagger
-keepattributes *Annotation*

# Retrofit & OkHttp
-keepattributes Signature, InnerClasses, EnclosingMethod
-dontwarn okio.**
-keep class com.aarcsx.krisho.core.network.dto.** { *; }

# Room
-keep class * extends androidx.room.RoomDatabase
-keep class com.aarcsx.krisho.core.local.room.** { *; }

# Razorpay
-keep class com.razorpay.** { *; }
-dontwarn com.razorpay.**

# Firebase
-keep class com.google.firebase.** { *; }

# Preserve line numbers for Crashlytics
-keepattributes SourceFile,LineNumberTable
-renamesourcefileattribute SourceFile
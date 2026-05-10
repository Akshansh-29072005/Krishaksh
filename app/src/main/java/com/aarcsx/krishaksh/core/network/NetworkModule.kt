package com.aarcsx.krishaksh.core.network

import com.aarcsx.krishaksh.BuildConfig
import com.aarcsx.krishaksh.core.network.api.*
import com.aarcsx.krishaksh.core.network.interceptor.AuthInterceptor
import com.aarcsx.krishaksh.core.network.interceptor.RefreshTokenAuthenticator
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import java.util.concurrent.TimeUnit
import javax.inject.Named
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
object NetworkModule {

    @Provides
    @Singleton
    @Named("base_url")
    fun provideBaseUrl(): String = BuildConfig.BASE_URL

    @Provides
    @Singleton
    fun provideLoggingInterceptor(): HttpLoggingInterceptor =
        HttpLoggingInterceptor().apply { 
            level = if (BuildConfig.DEBUG) {
                HttpLoggingInterceptor.Level.BODY 
            } else {
                HttpLoggingInterceptor.Level.NONE
            }
        }

    @Provides
    @Singleton
    @Named("authless")
    fun provideAuthlessOkHttp(logging: HttpLoggingInterceptor): OkHttpClient =
        OkHttpClient.Builder()
            .addInterceptor(logging)
            .connectTimeout(30, TimeUnit.SECONDS)
            .readTimeout(60, TimeUnit.SECONDS)
            .writeTimeout(60, TimeUnit.SECONDS)
            .build()

    @Provides
    @Singleton
    @Named("authless_retrofit")
    fun provideAuthlessRetrofit(@Named("base_url") baseUrl: String, @Named("authless") okHttpClient: OkHttpClient): Retrofit =
        Retrofit.Builder().baseUrl(baseUrl).client(okHttpClient).addConverterFactory(GsonConverterFactory.create()).build()

    @Provides
    @Singleton
    fun provideAuthApi(@Named("authless_retrofit") retrofit: Retrofit): AuthApiService = retrofit.create(AuthApiService::class.java)

    @Provides
    @Singleton
    fun provideOkHttp(
        authInterceptor: AuthInterceptor,
        authenticator: RefreshTokenAuthenticator,
        logging: HttpLoggingInterceptor
    ): OkHttpClient =
        OkHttpClient.Builder()
            .addInterceptor(authInterceptor)
            .authenticator(authenticator)
            .addInterceptor(logging)
            .connectTimeout(30, TimeUnit.SECONDS)
            .readTimeout(60, TimeUnit.SECONDS)
            .writeTimeout(60, TimeUnit.SECONDS)
            .build()

    @Provides
    @Singleton
    fun provideRetrofit(@Named("base_url") baseUrl: String, okHttpClient: OkHttpClient): Retrofit =
        Retrofit.Builder().baseUrl(baseUrl).client(okHttpClient).addConverterFactory(GsonConverterFactory.create()).build()

    @Provides @Singleton fun provideScanApi(retrofit: Retrofit): ScanApiService = retrofit.create(ScanApiService::class.java)
    @Provides @Singleton fun provideProductApi(retrofit: Retrofit): ProductApiService = retrofit.create(ProductApiService::class.java)
    @Provides @Singleton fun provideCommerceApi(retrofit: Retrofit): CommerceApiService = retrofit.create(CommerceApiService::class.java)
    @Provides @Singleton fun provideUserApi(retrofit: Retrofit): UserApiService = retrofit.create(UserApiService::class.java)
}

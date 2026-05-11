package com.aarcsx.krisho.core.auth

import android.content.Context
import com.google.android.gms.auth.api.signin.GoogleSignIn
import com.google.android.gms.auth.api.signin.GoogleSignInClient
import com.google.android.gms.auth.api.signin.GoogleSignInOptions
import dagger.hilt.android.qualifiers.ApplicationContext
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class GoogleSignInManager @Inject constructor(
    @ApplicationContext private val context: Context
) {
    private val gso = GoogleSignInOptions.Builder(GoogleSignInOptions.DEFAULT_SIGN_IN)
        .requestIdToken("248483110576-913ukuatgdpnacge6lp598stqlqiman0.apps.googleusercontent.com")
        .requestEmail()
        .build()

    val client: GoogleSignInClient = GoogleSignIn.getClient(context, gso)
}

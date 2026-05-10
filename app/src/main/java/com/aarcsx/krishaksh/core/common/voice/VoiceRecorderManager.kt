package com.aarcsx.krishaksh.core.common.voice

import android.content.Context
import android.media.MediaRecorder
import java.io.File
import java.io.FileOutputStream

class VoiceRecorderManager(private val context: Context) {
    private var mediaRecorder: MediaRecorder? = null
    private var currentFile: File? = null

    fun startRecording(): File? {
        val file = File(context.cacheDir, "audio_record_${System.currentTimeMillis()}.m4a")
        currentFile = file

        mediaRecorder = MediaRecorder(context).apply {
            setAudioSource(MediaRecorder.AudioSource.MIC)
            setOutputFormat(MediaRecorder.OutputFormat.MPEG_4)
            setAudioEncoder(MediaRecorder.AudioEncoder.AAC)
            setAudioSamplingRate(44100)
            setAudioEncodingBitRate(64000)
            setOutputFile(file.absolutePath)
            prepare()
            start()
        }
        return file
    }

    fun stopRecording() {
        mediaRecorder?.apply {
            stop()
            release()
        }
        mediaRecorder = null
    }

    fun getAmplitude(): Float {
        return mediaRecorder?.maxAmplitude?.toFloat() ?: 0f
    }
}

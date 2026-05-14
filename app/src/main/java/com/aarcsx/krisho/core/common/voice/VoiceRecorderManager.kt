package com.aarcsx.krisho.core.common.voice

import android.content.Context
import android.media.MediaRecorder
import kotlin.math.abs
import java.io.File

class VoiceRecorderManager(private val context: Context) {
    private var mediaRecorder: MediaRecorder? = null
    private var currentFile: File? = null
    private var recordingStartTime: Long = 0
    private val MAX_DURATION_MS = 60_000L // 60 seconds
    private val SILENCE_THRESHOLD = 2000 // amplitude threshold for silence

    fun startRecording(): File? {
        val file = File(context.cacheDir, "audio_record_${System.currentTimeMillis()}.m4a")
        currentFile = file
        recordingStartTime = System.currentTimeMillis()

        mediaRecorder = MediaRecorder(context).apply {
            setAudioSource(MediaRecorder.AudioSource.MIC)
            setOutputFormat(MediaRecorder.OutputFormat.MPEG_4)
            setAudioEncoder(MediaRecorder.AudioEncoder.AAC)
            setAudioSamplingRate(16000) // reduced from 44100 for compression
            setAudioEncodingBitRate(32000) // reduced from 64000 for aggressive compression
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

    fun getRecordingDuration(): Long {
        return System.currentTimeMillis() - recordingStartTime
    }

    fun hasReachedMaxDuration(): Boolean {
        return getRecordingDuration() >= MAX_DURATION_MS
    }

    fun isSilent(): Boolean {
        val amplitude = getAmplitude()
        return amplitude < SILENCE_THRESHOLD
    }

    fun getTrimmedFile(originalFile: File): File? {
        // Silence trimming will be handled by backend for stateless architecture
        // Return the original file; let server handle trimming if needed
        return originalFile
    }
}

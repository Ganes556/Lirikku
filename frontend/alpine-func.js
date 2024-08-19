import { MediaRecorder, register } from 'extendable-media-recorder';
import { connect } from 'extendable-media-recorder-wav-encoder';
import htmx from 'htmx.org';
export function audioRecorder(maxRecordingTime) {
  return {
    isRecording: false,
    mediaRecorder: null,
    audioChunks: [],
    errorMessage: '',
    maxRecordingTime: maxRecordingTime, // seconds in milliseconds
    recordingTimeout: null,

    async startRecording() {
      this.errorMessage = '';

      try {
        await register(await connect());

        const stream = await navigator.mediaDevices.getUserMedia({
          audio: true,
        });

        // Initialize MediaRecorder with the stream and MIME type

        // const audioContext = new AudioContext({ sampleRate: 44100 });

        this.mediaRecorder = new MediaRecorder(stream, {
          mimeType: 'audio/wav',
        });

        // Start recording and update state
        this.mediaRecorder.start();
        this.isRecording = true;

        // Automatically stop recording after maxRecordingTime
        this.recordingTimeout = setTimeout(() => {
          this.stopRecording();
        }, this.maxRecordingTime);

        // Collect audio data as it's available
        this.mediaRecorder.addEventListener('dataavailable', (event) => {
          this.audioChunks.push(event.data);
        });

        // Handle the stop event to process the recording
        this.mediaRecorder.addEventListener('stop', async () => {
          const audioBlob = new Blob(this.audioChunks, {
            type: this.mediaRecorder.mimeType,
          });

          const form = document.getElementById('search-recorded');
          form.setAttribute('hx-encoding', 'multipart/form-data');
          form.setAttribute('hx-trigger', 'none');

          const fileBlob = new File([audioBlob], 'audio.wav');
          const formData = new FormData(form);
          formData.append('audio', fileBlob);
          htmx.ajax('POST', form.getAttribute('hx-post'), {
            values: {
              audio: formData.get('audio'),
            },
            source: form,
          });

          this.isRecording = false;
          this.audioChunks = [];
          clearTimeout(this.recordingTimeout); // Clear the timeout if stop is triggered manually
        });
      } catch (error) {
        this.errorMessage = 'Error accessing microphone: ' + error.message;
      }
    },

    stopRecording() {
      if (this.mediaRecorder && this.isRecording) {
        this.mediaRecorder.stop();
      }
      this.isRecording = false;
    },
  };
}

// hello i have a question, why the htmx ajax not send the request properly?,  this is my code,

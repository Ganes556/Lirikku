import { MediaRecorder, register } from 'extendable-media-recorder';
import { connect } from 'extendable-media-recorder-wav-encoder';
import htmx from 'htmx.org';

export function audioRecorder(maxRecordingTime) {
  return {
    isRecording: false,
    mediaRecorder: null,
    audioChunks: [],
    maxRecordingTime: maxRecordingTime,
    retries: 0,
    searchByRecorded: false,
    retriesQparam: {
      identifier: crypto.randomUUID(),
      samplems: maxRecordingTime,
    },
    recordingTimeout: null,
    stream: null,
    encoderRegistered: false, // New flag to track encoder registration

    async init() {
      try {
        if (!this.encoderRegistered) {
          await register(await connect());
          this.encoderRegistered = true; // Mark encoder as registered
        }
      } catch (error) {
        console.log('Error registering encoder: ' + error.message);
      }
    },

    async startRecording() {
      this.searchByRecorded = true;
      htmx.find('#search-results').innerHTML = '';
      htmx.find('#error-results').innerHTML = '';
      try {
        // Ensure previous stream is stopped
        if (this.stream) {
          this.stream.getTracks().forEach((track) => track.stop());
        }

        this.stream = await navigator.mediaDevices.getUserMedia({
          audio: true,
        });

        // Reinitialize MediaRecorder
        this.mediaRecorder = new MediaRecorder(this.stream, {
          mimeType: 'audio/wav',
        });

        this.isRecording = true;
        this.mediaRecorder.start();

        this.recordingTimeout = setTimeout(() => {
          this.stopRecording();
        }, this.maxRecordingTime);

        this.mediaRecorder.addEventListener('dataavailable', (event) => {
          this.audioChunks.push(event.data);
        });

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

          // const base64res = await this.blobToBase64(fileBlob)
          // :hx-on::after-request="afterRequest()"
          // console.log(base64res)
          // console.log('param -> ', this.retriesQparam);

          htmx.ajax('POST', form.getAttribute('hx-post'), {
            values: {
              ...this.retriesQparam,
              audio: formData.get('audio'),
            },
            source: form,
          });
        });
      } catch (error) {
        console.log('Error accessing microphone: ' + error.message);
      }
    },
    async stopRecording() {
      if (this.mediaRecorder && this.isRecording) {
        this.isRecording = false;
        this.audioChunks = [];
        clearTimeout(this.recordingTimeout);
        this.mediaRecorder.stop();
        this.stream.getTracks().forEach((track) => track.stop()); // Stop the stream
      }
    },
    // afterRequest() {
    //   htmx.on('htmx:afterRequest', (e) => {
    //     if (this.searchByRecorded && e.detail.successful) {
    //       if (e.detail.xhr.status != 200) {
    //         return;
    //       }
    //       if (this.retries < 2) {
    //         this.retries += 1;
    //         this.retriesQparam.samplems += this.maxRecordingTime;
    //         this.startRecording();
    //       } else {
    //         htmx.find('#error-results').innerHTML = "couldn't quite catch that";
    //         this.resetRetries();
    //       }
    //     }
    //   });
    // },
    resetRetries() {
      this.retries = 0;
      this.searchByRecorded = false;
      this.retriesQparam.identifier = crypto.randomUUID();
      this.retriesQparam.samplems = this.maxRecordingTime;
    },
  };
}

// function blobToBase64(blob) {
//   return new Promise((resolve, reject) => {
//     const reader = new FileReader();
//     reader.onloadend = () => resolve(reader.result.split(',')[1]);
//     reader.onerror = reject;
//     reader.readAsDataURL(blob);
//   });
// }

// await disconnect(); // Disconnect the encoder

// type ReqSearchAudio struct {
// 	Identifier string `query:"identifier"`
// 	Timestamp int64 `query:"timestamp"`
// 	Samplems int `query:"samplems"`
// }

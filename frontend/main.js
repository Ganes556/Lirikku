import './style.css';
import 'htmx.org';
import './htmx-ext/json-enc.js';
import './htmx-ext/hx-target.js';
import './node_modules/@material-tailwind/html/scripts/dialog.js';
import './node_modules/@material-tailwind/html/scripts/ripple.js';
import { audioRecorder } from './alpine-func.js';
import Alpine from 'alpinejs';

window.htmx = htmx;
window.Alpine = Alpine;

// window.onload = function () {  
//     htmx.config.responseTargetUnsetsError = false; htmx.config.responseTargetSetsError = true;
// }

Alpine.data('audioRecorder', audioRecorder)
Alpine.start();
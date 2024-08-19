import './style.css';
import 'htmx.org';
import './htmx-ext/json-enc.js';
import './node_modules/@material-tailwind/html/scripts/dialog.js';
import './node_modules/@material-tailwind/html/scripts/ripple.js';
import { audioRecorder } from './alpine-func.js';
import Alpine from 'alpinejs';

window.htmx = htmx;
window.Alpine = Alpine;

Alpine.data('audioRecorder', audioRecorder)
Alpine.start();
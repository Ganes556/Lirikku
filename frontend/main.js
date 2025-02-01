import './style.css';
import 'htmx.org';
import './htmx-ext/json-enc.js';
import './htmx-ext/hx-target.js';
import './node_modules/@material-tailwind/html/scripts/dialog.js';
import './node_modules/@material-tailwind/html/scripts/ripple.js';
import Clipboard from '@ryangjchandler/alpine-clipboard';
import { audioRecorder } from './alpine-func.js';
import Alpine from 'alpinejs';
import AsyncAlpine from 'async-alpine';

Alpine.plugin(Clipboard);
window.htmx = htmx;
window.Alpine = Alpine;

Alpine.data('audioRecorder', audioRecorder);
Alpine.start();

// Vue.js
import { initVue } from './vue.ts';

// TailwindCSS Styles
import "../src/tailwind.css";

// Custom font
// import "./fonts.css";

// Landing page 
import { env } from './config.ts';

import ApexCharts from 'apexcharts'

declare global {
    interface Window {
        vm: any, // Vue.js instance
        observer: any
    }
}

document.addEventListener('DOMContentLoaded', () => {
    
    initVue();

});

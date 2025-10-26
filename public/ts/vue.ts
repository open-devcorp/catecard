import { createApp } from 'vue/dist/vue.esm-bundler.js';
import { env } from './config';

import TestComponent from './components/test-component.vue';
import ModalComponent from './components/modal.component.vue';


import { ref } from 'vue';

function NewApp() {
    const app = createApp({
        data() {
            return {
                showModal: false,
                isEditMode: false,
            }
        }
    });
    app.component('test-component', TestComponent);
    app.component('modal-component', ModalComponent);
    return app;
}

let lastApp: any = null;

export async function initVue(callbackFn?: Function) {
    if (window.vm) {
        lastApp?.unmount()
        lastApp = null;
    }
    let mountPoint = document.querySelector("[data-behavior='vue']");
    if (!mountPoint) {
        if (env == "dev") {
            console.warn("vue mount not found")
        }
        return;
    }

    // Save a copy of the original HTML
    const originalHTML = mountPoint.innerHTML;

    // Extract <script> and <style> tags so Vue's template compiler doesn't ignore them
    const extractedScripts: Array<{
        src?: string | null,
        type?: string | null,
        async?: boolean,
        defer?: boolean,
        content?: string | null,
        attrs: Array<{ name: string, value: string }>
    }> = [];
    const extractedStyles: HTMLElement[] = [];

    try {
        const scripts = Array.from(mountPoint.querySelectorAll('script')) as HTMLScriptElement[];
        for (const s of scripts) {
            extractedScripts.push({
                src: s.getAttribute('src'),
                type: s.getAttribute('type'),
                async: s.hasAttribute('async'),
                defer: s.hasAttribute('defer'),
                content: s.src ? null : s.innerHTML,
                attrs: Array.from(s.attributes).map(a => ({ name: a.name, value: a.value }))
            });
            // remove from DOM so Vue's compiler doesn't see it
            s.parentElement?.removeChild(s);
        }

        const styles = Array.from(mountPoint.querySelectorAll('style')) as HTMLElement[];
        for (const st of styles) {
            // clone and store; we'll append to head after mount
            extractedStyles.push(st.cloneNode(true) as HTMLElement);
            st.parentElement?.removeChild(st);
        }
    } catch (e) {
        // ignore extraction errors and continue
        console.warn('Error extracting scripts/styles before Vue mount', e);
    }



    let app = NewApp();
    try {
        window.vm = app.mount("[data-behavior='vue']")
        if (callbackFn) callbackFn()
        // After mount, re-insert styles and re-run scripts in original order
        try {
            // append styles to head
            if (extractedStyles.length > 0) {
                for (const st of extractedStyles) {
                    try { document.head.appendChild(st); } catch (e) { console.warn('Failed to append style', e); }
                }
            }

            // re-run scripts in sequence to preserve order
            for (const desc of extractedScripts) {
                try {
                    // If this script had a src and it already exists in the document, skip to avoid duplicates
                    if (desc.src) {
                        const exists = !!document.querySelector(`script[src="${desc.src}"]`);
                        if (exists) {
                            // Wait a tick to preserve ordering
                            await Promise.resolve();
                            continue;
                        }
                    }

                    const el = document.createElement('script');
                    for (const a of desc.attrs) {
                        try { el.setAttribute(a.name, a.value); } catch (e) { /* ignore */ }
                    }

                    if (desc.src) {
                        el.src = desc.src;
                        // ensure scripts load and execute; append to body
                        document.body.appendChild(el);
                        // wait for it to load before continuing to next script
                        await new Promise<void>((res) => {
                            el.onload = () => res();
                            el.onerror = () => { console.warn('Script load error', desc.src); res(); };
                        });
                    } else if (desc.content) {
                        // Wrap inline script content in an IIFE to avoid redeclaring const/let in global scope
                        try {
                            el.text = `(function(){\n${desc.content}\n})();`;
                        } catch (e) {
                            el.text = desc.content;
                        }
                        document.body.appendChild(el);
                    } else {
                        document.body.appendChild(el);
                    }
                } catch (se) {
                    console.warn('Error reinserting script', se);
                }
            }
        } catch (re) {
            console.warn('Error reinserting extracted assets', re);
        }
    } catch (e) {
        alert("Vue mount error:" + e);
        console.error("Vue mount error", e);
        // Restore the original HTML
        mountPoint.innerHTML = originalHTML;
    }
    lastApp = app;
}

import {createApp} from 'vue'
import App from './App.vue'
import "bootstrap/dist/css/bootstrap.min.css"
import 'bootstrap-icons/font/bootstrap-icons.css'
import "bootstrap"
import Vue3Toasity, { type ToastContainerOptions } from 'vue3-toastify';
import 'vue3-toastify/dist/index.css';
import JsonEditorVue from "json-editor-vue";
import router from './router'

createApp(App).use(
    JsonEditorVue, {
        mainMenuBar: false,
        navigationBar: false,
        statusBar: false,
        mode: 'text',
        queryLanguagesIds: ['javascript', 'lodash', 'jmeshpath']
}).use(
    Vue3Toasity,
    {
        autoClose: 3000,
    } as ToastContainerOptions,
)
.mount('#app')

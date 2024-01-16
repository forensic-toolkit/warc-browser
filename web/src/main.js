import { createApp } from 'vue'
import App from './App.vue'
import WaveUI from 'wave-ui'
import 'wave-ui/dist/wave-ui.css'
import '@mdi/font/css/materialdesignicons.min.css'

createApp(App).
    use(WaveUI, {
        /* Wave UI options */
        on: "#app",
        theme: "dark",
        // breakpoint: "xl",
        css: { grid: 12 },
        colors: {
            dark:{
                // primary: '#9ac332',
                // secondary: '#5d9a26',
                // // Custom color names should be kebab-case.
                // 'mint-green': '#bff8db',
            },
        },
        presets: {
            "replay-web-page": {
                height: "80vh",
                width: "100%",
            },
            "w-tree": {"height": "100%"},
        }
    }).
    mount('#app')
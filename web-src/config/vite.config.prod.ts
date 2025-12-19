import { mergeConfig,loadEnv } from 'vite';
import baseConig from './vite.config.base';
import configCompressPlugin from './plugin/compress';
import configVisualizerPlugin from './plugin/visualizer';

const BASE_URL = loadEnv('production', process.cwd()).VITE_BASE_URL; 

export default mergeConfig(
  baseConig,
  {
    mode: 'production',
    base: BASE_URL,
    mock: true,
    plugins: [configCompressPlugin('gzip'), configVisualizerPlugin()],
    build: {
      rollupOptions: {
        output: {
          manualChunks: {
            // vcharts: ['vue-echarts'],
            // tinymce: ['tinymce'],
            // tinymcevue: ['@tinymce/tinymce-vue'],
            // lodash: ['lodash'],
            echarts: ['echarts'],
            echarts4: ['echarts4'],
            qs: ['query-string'],
            // vaceEditor: ['vue3-ace-editor'],
            // clipboard3: ['vue-clipboard3'],
            opentinyIcon: ['@opentiny/vue-icon'],
            pinia: ['pinia'],
            axios: ['axios'],
            dayjs: ['dayjs'],
            mitt: ['mitt'],
            vue: ['vue', 'vue-router','@vueuse/core', 'vue-i18n'],
          },
          // 文件加上时间戳
          // chunkFileNames: `js/[name].[hash].js`,
          // entryFileNames: `js/[name].[hash].js`,
          // assetFileNames: `[ext]/[name].[hash].[ext]`,
        },
      },
      chunkSizeWarningLimit: 500,
    },
  }
  
);

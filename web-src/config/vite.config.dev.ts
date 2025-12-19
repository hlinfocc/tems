import { mergeConfig, loadEnv } from 'vite';
import eslint from 'vite-plugin-eslint';
import baseConfig from './vite.config.base';
import { createProxy } from './plugin/proxy';

// const proxyConfig = ()=>{
//   let proxtData = loadEnv('development', process.cwd()).VITE_API_URL_PROXY_PREFIX;
//   console.log("proxtData",proxtData);
//   let ret = createProxy(proxtData);
//   return ret;
// };

const proxyConfig = {
  [loadEnv('development', process.cwd()).VITE_BASE_API]: {
    target: loadEnv('development', process.cwd()).VITE_SERVER_HOST,
    changeOrigin: true,
    logLevel: 'debug',
    rewrite: (path) =>
      path.replace(
        new RegExp(`${loadEnv('development', process.cwd()).VITE_BASE_API}`),
        ''
      ),
  },
};

export default mergeConfig(
  {
    mode: 'development',
    server: {
      open: true,
      host: '0.0.0.0',
      fs: {
        strict: true,
      },
      // proxy: proxyConfig,
      proxy: {
        ...proxyConfig,
      },
    },
    plugins: [
      eslint({
        include: ['src/**/*.ts', 'src/**/*.tsx', 'src/**/*.vue'],
        exclude: ['node_modules'],
      }),
    ],
  },
  baseConfig
);

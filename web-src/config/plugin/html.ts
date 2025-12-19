/**
 * Plugin to minimize and use ejs template syntax in index.html.
 * https://github.com/anncwb/vite-plugin-html
 */
import type { PluginOption } from 'vite';
import { createHtmlPlugin } from 'vite-plugin-html';
import pkg from '../../package.json';
// import { GLOB_CONFIG_FILE_NAME } from '../constant';
// import { loadEnv } from 'vite';
// import type { ConfigEnv } from 'vite';


export default function configHtmlPlugin(title,publicPath) {
  console.log("process.env.:",process.env.NODE_ENV);
  // const root = process.cwd();
  // const env = loadEnv(nodeEnv, root);
  
  // const title = env.VITE_GLOB_APP_TITLE;
  // const publicPath = env.VITE_PUBLIC_PATH;

  const isBuild = true;
  const isAppConfig = false;
  
  const path = publicPath.endsWith('/') ? publicPath : `${publicPath}/`;
  const GLOB_CONFIG_FILE_NAME = '_app.config.js';
  
  const getAppConfigSrc = () => {
    return `${path || '/'}${GLOB_CONFIG_FILE_NAME}?v=${pkg.version}-${new Date().getTime()}`;
  };

  const htmlPlugin: PluginOption[] = createHtmlPlugin({
    minify: isBuild,
    inject: {
      // Inject data into ejs template
      data: {
        title: title,
      },
      // Embed the generated app.config.js file
      tags: isAppConfig
        ? [
            {
              tag: 'script',
              attrs: {
                src: getAppConfigSrc(),
              },
            },
          ]
        : [],
    },
  });
  return htmlPlugin;
}
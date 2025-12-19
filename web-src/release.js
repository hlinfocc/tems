const { log } = require('console');
const fs = require('fs');
const path = require('path');

const distPath = path.join(__dirname, '../web-api/assets/dist');

if(fs.existsSync(distPath)){
    fs.rmdirSync(distPath, { recursive: true });

    if(!fs.existsSync(distPath)){
        console.log('成功删除文件夹: ../web-api/assets/dist');
    }
}

const sourcePath = path.join(__dirname, './dist');

if(fs.existsSync(sourcePath)){
    fs.rename(sourcePath, distPath, (err) => {
      if (err) {
        console.error('移动文件夹失败:', err);
      } else {
        console.log('成功移动文件夹:', sourcePath, '到', distPath);
      }
    });
}

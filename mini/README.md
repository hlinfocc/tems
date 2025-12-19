# 评教系统

# 功能介绍

# 主要文件：

- 基础配置文件 ：
  - `app.json` ：定义了全局配置、页面路径和底部tab栏
  - `app.js` ：设置了小程序全局逻辑和数据管理
  - `app.wxss` ：定义了全局样式，使用蓝色(#1989fa)作为主色调
  - `project.config.json` ：微信小程序项目配置文件
页面文件：

- 首页 ：
  
  - `index.wxml` ：包含轮播图和待评教课程列表
  - `index.wxss` ：首页样式设计
  - `index.js` ：首页数据加载和交互逻辑
- 我的页面 ：
  
  - `profile.wxml` ：展示用户信息和功能菜单
  - `profile.wxss` ：个人页面样式
  - `profile.js` ：个人页面逻辑处理
- 评教页面 ：
  
  - `evaluation.wxml` ：评教表单界面
  - `evaluation.wxss` ：评教页面样式
  - `evaluation.js` ：评教功能的交互逻辑
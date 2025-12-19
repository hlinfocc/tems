# 使用方法

* AceEditor：在线代码编辑器
* TinymceEditor：富文本编辑器

```html
<template>
  <div>
  <button @click="gotopy">更换Python</button>
  <button @click="getval">获取值</button>
  <ace-editor v-model:value="content" :lang="lang" :height="height"/>

  <cy-fluent-editor v-model:value="content" />
  </div>
</template>
<script setup lang="ts">
import { ref ,watch} from "vue";
import { AceEditor,CyFluentEditor } from '@/components/editor';

const height = ref('300px');
//语言： 0:c/c++ ,6:python
const lang = ref<number>(0);
const content = ref<string>('');


function gotopy(){
  lang.value = 6;
  console.log('gotopy:',lang.value);
}
function getval(){
  console.log('content:',content.value);
}
</script>
```
<template>
  <div :style="`height:${height};width:${width};`">
    <v-ace-editor
    v-model:value="content"
    :lang="aceLang"
    :options="aceOptions"
    :theme="theme"
    :readonly="readOnly"
    :style="`height:${height};width:100%;`" />
  </div>
</template>

<script setup lang="ts">
import {onMounted, ref, defineProps, watch,toRef, toRaw } from "vue";
// @ts-ignore
import { VAceEditor } from "vue3-ace-editor"
import './ace-config';

const props = defineProps(["value","lang","height","width","readonly"]);
// 将props转为响应式
const propsValue = toRef(props, 'value')
const propsLang = toRef(props, 'lang')
const propsHeight = toRef(props, 'height')
const propsWidth = toRef(props, 'width')
const propsReadonly = toRef(props, 'readonly')
const content = ref<string>('')
const aceLang = ref<string>('html')
const height = ref<string>('300px')
const width = ref<string>('100%')
const theme = ref<string>('eclipse')
const readOnly = ref<boolean>(false)

// Emits声明
const emit = defineEmits(['update:value','change']);

const aceOptions=ref({
    enableBasicAutocompletion: true,
    enableSnippets: true,
    enableLiveAutocompletion: true,
    fontSize: 14,
    tabSize: 2,
    showPrintMargin: false,
    highlightActiveLine: true,
    useWorker: true,
    wrap:'free'
})

watch(propsValue, (newValue, oldValue) => {
    content.value = propsValue.value
});
watch(propsLang, (newValue, oldValue) => {
    initLang();
});
watch(propsHeight, (newValue, oldValue) => {
    initHeight();
});
watch(propsWidth, (newValue, oldValue) => {
    initHeight();
});
watch(propsReadonly, (newValue, oldValue) => {
    initHeight();
});
// 监听值变化，将值传回父组件
watch(content, (newValue, oldValue) => {
    emit('update:value',newValue);
    emit('change',newValue);
});

onMounted(()=>{
  initHeight();
  initLang();
  content.value = propsValue.value
})

function initLang(){
  if(propsLang.value === 1){
    aceLang.value = 'c_cpp';
  }else if(propsLang.value === 6){
    aceLang.value = 'python';
  }else if(propsLang.value === 43){
    aceLang.value = 'text';
  }else if(propsLang.value === 7){
    aceLang.value = 'sql';
  }else if(propsLang.value === 8){
    aceLang.value = 'javascript';
  }else if(propsLang.value === 9){
    aceLang.value = 'json';
  }else{
    aceLang.value = 'html';
  }
}

function initHeight(){
  if(propsHeight.value){
    height.value = propsHeight.value
  }
  if(propsWidth.value){
    width.value = propsWidth.value
  }
  if(propsReadonly.value){
    readOnly.value = propsReadonly.value
  }
}

</script>
<script lang="ts">
export default {
  name: "AceEditor"
};
</script>
<style lang="less"  scoped>

.row-container:not(:last-child) {
  margin-bottom: 16px;
}
</style>

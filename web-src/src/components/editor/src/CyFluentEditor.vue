<template>
    <div>
          <div v-loading="loading" size="large" tiny-loading__text="loadingTxt">
            <tiny-fluent-editor :id="editorId" v-model="myContentValue" :options="options"></tiny-fluent-editor>
          </div>
    </div>
  </template>
  <script lang="ts" setup>
    import { ref, onMounted, toRaw,toRef,watch, toRefs,defineProps } from 'vue';
    import { Loading,Modal,TinyModal } from '@opentiny/vue'
    import { TinyFluentEditor } from '@opentiny/vue'
    import type { I18N } from '@opentiny/fluent-editor'

    // @ts-ignore
    import { getToken } from '@/utils/auth';

    const vLoading = Loading.directiv;

    const loading = ref(false);
    const loadingTxt = ref('');
  
    const myContentValue = ref('');
    const editorId = ref('fluentEditor');
  
    const props = defineProps({
            id: {
                type: String,
                default: 'fluentEditor'
            },
            value: {
                type: String,
                default: ''
            },
            menubar: {
                type: Boolean,
                default: true
            },
            readOnly: {
                type: Boolean,
                default: true // 是否禁用编辑
            },
            // 基本路径，默认为空根目录，如果你的项目发布后的地址为目录形式，
            baseUrl: {
                type: String,
                default: './'
            },
            disabled: {
                type: Boolean,
                default: false
            },
            
    })
    const options = ref({
    placeholder: '请输入内容',
    modules: {
        // 工具栏
        toolbar: [
            ['undo', 'redo', 'clean', 'format-painter'],
            [
                // 请保留默认值为 false
                { header: [1, 2, 3, 4, 5, 6, false] },
                { font: [false, '仿宋_GB2312, 仿宋', '楷体', '隶书', '黑体', '无效字体, 隶书'] },
                { size: [false, '12px', '14px', '16px', '18px', '20px', '24px', '32px', '36px', '48px', '72px'] },
                { 'line-height': [false, '1.2', '1.5', '1.75', '2', '3', '4', '5'] },
            ],
            ['bold', 'italic', 'strike', 'underline', 'divider'],
            [{ color: [] }, { background: [] }],
            [{ align: '' }, { align: 'center' }, { align: 'right' }, { align: 'justify' }],
            [{ list: 'ordered' }, { list: 'bullet' }, { list: 'check' }],
            [{ script: 'sub' }, { script: 'super' }],
            [{ indent: '-1' }, { indent: '+1' }],
            [{ direction: 'rtl' }],
            ['link', 'blockquote', 'code', 'code-block'],
            ['image', 'file'],
            ['emoji', 'video', 'formula', 'screenshot', 'fullscreen'],
            ['better-table'],
        ],
        // 字数统计
        counter: {
            count:1000000
        },
        i18n: {
          lang: 'zh-CN',
        },
        syntax: false,
        mathlive: true,
        // 表格
        'better-table': {
            operationMenu: {
                color: {
                    text: '主题色',
                    colors: [
                        '#ffffff',
                        '#f2f2f2',
                        '#dddddd',
                        '#a6a6a6',
                        '#666666',
                        '#000000',
                        '#c00000',
                        '#ff0000',
                        '#ffc8d3',
                        '#ffc000',
                        '#ffff00',
                        '#fff4cb',
                        '#92d050',
                        '#00b050',
                        '#dff3d2',
                        '#00b0f0',
                        '#0070c0',
                        '#d4f1f5',
                        '#002060',
                        '#7030a0',
                        '#7b69ee',
                        '#1476ff',
                        '#ec66ab',
                        '#42b883'
                        ]
                    }
                }
            }
        }
    })
  const { id } = toRefs(props);
  const propsValue = toRef(props, 'value')
  watch(propsValue, (newValue, oldValue) => {
    myContentValue.value = propsValue.value
  });
  
  // Emits声明
  const emit = defineEmits(['update:value','onClick','onChange']);
  // const files = ref([]);
  const replaceImgSrcPrefixByUrl = (url:any)=>{return url;};
  
  
  onMounted(()=>{
    myContentValue.value = propsValue.value
      
  });
  // 监听值变化，将值传回父组件
  watch(myContentValue, (newValue, oldValue) => {
      emit('update:value',myContentValue.value);
  });
  
  /* function handlePaiBan() {
      if(!myValue.value || myValue.value == ''){
          Modal.message({ message: '请先填写内容', status: 'warning' })
          return false;
      }
      loading.value = true;
      loadingTxt.value = "正在处理中，请稍后。。。";
      // @ts-ignore
      typesetingApi({html:myValue.value}).then( h => {
          // @ts-ignore
          if (h.code == 200){
              myValue.value = h.data;
              onChange(h.data);
              Modal.message({ message: '处理成功！', status: 'success' });
              loading.value = false;
          }else{
              // @ts-ignore
              Modal.message({ message: h.msg, status: 'error' });
          }
          loading.value = false;
      }).catch(()=>{
          loading.value = false;
      })
  } */
  
   
    // 设置内容
    function setContent(content) {
      
    }
    // 在当前光标处插入内容
    function insertContent(content) {
      
    }
    // 数据清狂
    function clear() {
        myContentValue.value = ''
    }
   
  </script>
  <script lang="ts">
  export default {
    name: "CyFluentEditor"
  };
  </script>
  <style lang="less"  scoped>
  @import '@opentiny/fluent-editor/style.css';

  :global(.tox-promotion){
    display: none;
  }
  :global(.tox .tox-mbtn){
    width: 50px !important;
  }
  :global(.tox .tox-textarea-wrap){
    height: 100% !important;
  }
  :global(.tox .tox-textarea){
    height: 100% !important;
  }
  </style>
  
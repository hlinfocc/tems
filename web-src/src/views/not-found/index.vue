<template>
  <div class="container">
    <div class="content">
      <div class="content-main">
        <img src="@/assets/images/500.png" alt="404" class="image" />
        <h3 class="tip">页面未找到</h3>
        <span>将{{ currTime }}秒后自动跳转到首页</span>
      </div>
    </div>
  </div>
</template>



<script lang="ts" setup>
  import { useRouter } from 'vue-router';
  import {ref, reactive, toRefs, onMounted, toRaw} from 'vue';

  const router = useRouter();
  const back = () => {
    router.push({ name: 'Home' });
  };
  const currTime = ref<number>(6);
  let timerInterval:any;
  const timeLoop = ()=>{
    if(currTime.value===1){
      clearInterval(timerInterval);
      back();
      return;
    }
    currTime.value -= 1;
  }
  onMounted(()=>{
    timerInterval = setInterval(()=>{
      timeLoop();
    },1000)
  })
</script>


<style scoped lang="less">
  @import '@/assets/style/exception.less'; /* 引入公共样式 */
</style>

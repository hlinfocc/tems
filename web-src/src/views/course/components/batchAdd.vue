<template>
  <tiny-drawer
    v-model:visible="visible"
    :title="'批量新增课程'"
    :width="`40%`"
    show-footer
    :mask-closable="false"
    @close="handleClose"
    @confirm="handleSubmit"
  >
    <tiny-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-position="top"
      size="large"
    >
      <tiny-row :gutter="20">
        <tiny-col :span="24">
          <tiny-form-item label="课程名称" prop="courseName" extra="输入课程名称，每行一个">
            <tiny-input
            type="textarea"
            rows="20"
            v-model="formData.courseName"
            placeholder="请输入课程名称（每行一个）"
          />
          </tiny-form-item>
        </tiny-col>
      </tiny-row>
    </tiny-form>
    <template #footer>
      <tiny-button @click="handleClose">取消</tiny-button>
      <tiny-button type="info" @click="handleSubmit" :loading="loading">
        保存
      </tiny-button>
    </template>
  </tiny-drawer>
</template>

<script setup lang="ts">
  import {
    ref,
    reactive,
    watch,
    toRaw,
    toRef,
    onMounted,
    defineEmits,
    defineProps,
  } from 'vue';
  import {
    Drawer as TinyDrawer,
    Form as TinyForm,
    FormItem as TinyFormItem,
    Input as TinyInput,
    Button as TinyButton,
    Row as TinyRow,
    Col as TinyCol,
    RadioGroup as TinyRadioGroup,
    Radio as TinyRadio,
    Modal,
  } from '@opentiny/vue';
  import { batchAddCourse, getCourseDetail } from '@/api/course';
  import type { CourseInfo, CourseForm } from '@/api/course';

  const visible = ref<boolean>(false);
  const props = defineProps<{
    visible: boolean;
    courseData?: CourseInfo;
  }>();

  const emit = defineEmits(['update:visible', 'success']);
  const propsVisible = toRef(props, 'visible');
  const propsCourseData = toRef(props, 'courseData');

  const formRef = ref();
  const isEdit = ref(false);
  const loading = ref(false);
  const formData = reactive<CourseForm>({
    id: undefined,
    courseName: '',
    status: 0,
  });

  const rules = {
    courseName: [
      {
        required: true,
        message: '请输入课程名称',
        trigger: 'blur',
      },
      {
        min: 1,
        max: 50,
        message: '课程名称长度在 1 到 50 个字符',
        trigger: 'blur',
      },
    ],
    status: [
      {
        required: true,
        message: '请选择状态',
        trigger: 'change',
      },
    ],
  };

  // 重置表单
  const resetForm = () => {
    formData.id = undefined;
    formData.courseName = '';
    formData.status = 0;
    formRef.value?.resetFields();
  };

  // 关闭抽屉
  const handleClose = () => {
    visible.value = false;
    emit('update:visible', false);
    setTimeout(() => {
      resetForm();
    }, 300);
  };

  // 提交表单
  const handleSubmit = async () => {
    if (!formRef.value) return;

    formRef.value?.validate(async (valid: boolean) => {
      if (valid) {
        try {
          loading.value = true;
          const submitData = { ...toRaw(formData) };
          if (!submitData.courseName) {
            return;
          }
          const courseList = submitData.courseName
            .split(/\r?\n/)
            .map((line) => {
              return {
                courseName: line.trim(),
                status: 0,
              };
            });
          let res;

          res = await batchAddCourse(courseList);

          if (res.code === 200) {
            Modal.message({
              message: res.msg,
              status: 'success',
            });
            emit('success');
            handleClose();
          } else {
            Modal.message({
              message: res.msg,
              status: 'error',
            });
          }
        } catch (error: any) {
          Modal.message({
            message: error.message || '操作失败',
            status: 'error',
          });
        } finally {
          loading.value = false;
        }
      }
    });
  };

  
  watch(visible, (newValue) => {
    emit('update:visible', newValue);
  });

  watch(propsVisible, (newValue) => {
    visible.value = newValue;
  });

  onMounted(() => {
    visible.value = propsVisible.value;
  });
</script>

<style scoped>
  .btn-group {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }

  .btn-group :deep(.tiny-btn) {
    margin-left: 10px;
  }
</style>

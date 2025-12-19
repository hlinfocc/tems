<template>
  <tiny-drawer
    v-model:visible="visible"
    :title="isEdit ? '编辑课程' : '新增课程'"
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
          <tiny-form-item label="课程名称" prop="courseName">
            <tiny-input
              v-model="formData.courseName"
              placeholder="请输入课程名称"
              :maxlength="50"
              show-word-limit
            ></tiny-input>
          </tiny-form-item>
        </tiny-col>
        <tiny-col :span="24">
          <tiny-form-item label="状态" prop="status">
            <tiny-radio-group v-model="formData.status">
              <tiny-radio :label="0">正常</tiny-radio>
              <tiny-radio :label="1">禁用</tiny-radio>
            </tiny-radio-group>
          </tiny-form-item>
        </tiny-col>
      </tiny-row>
    </tiny-form>
    <template #footer>
      <tiny-button @click="handleClose">取消</tiny-button>
      <tiny-button type="info" @click="handleSubmit" :loading="loading">
        {{ isEdit ? '更新' : '添加' }}
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
  import { addCourse, updateCourse, getCourseDetail } from '@/api/course';
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
          let res;

          if (isEdit.value) {
            res = await updateCourse(submitData);
          } else {
            res = await addCourse(submitData);
          }

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

  // 监听数据变化
  watch(
    () => propsCourseData.value,
    async (val) => {
      if (val && val.id) {
        isEdit.value = true;
        try {
          loading.value = true;
          const res = await getCourseDetail(val.id);
          if (res.code === 200) {
            Object.assign(formData, res.data);
          } else {
            Modal.message({
              message: res.msg || '获取课程详情失败',
              status: 'error',
            });
            handleClose();
          }
        } catch (error) {
          Modal.message({
            message: '获取课程详情失败',
            status: 'error',
          });
          handleClose();
        } finally {
          loading.value = false;
        }
      } else {
        isEdit.value = false;
        resetForm();
      }
    },
    { immediate: true, deep: true },
  );

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

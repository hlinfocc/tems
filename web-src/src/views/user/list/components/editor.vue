<template>
  <tiny-drawer
    v-model:visible="visible"
    :title="isEdit ? '编辑管理员' : '新增管理员'"
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
      :hide-required-asterisk="false"
      label-position="top"
      size="large"
    >
      <tiny-row :gutter="20">
        <tiny-col :span="12">
          <tiny-form-item label="姓名" prop="name">
            <tiny-input
              v-model="formData.name"
              placeholder="请输入姓名"
              :maxlength="100"
              show-word-limit
              @blur="handleNameInput"
            ></tiny-input>
          </tiny-form-item>
        </tiny-col>
        <tiny-col :span="12">
          <tiny-form-item label="用户名" prop="username">
            <tiny-input
              v-model="formData.username"
              placeholder="请输入用户名"
              :maxlength="50"
              show-word-limit
            ></tiny-input>
          </tiny-form-item>
        </tiny-col>
        <tiny-col :span="12" v-if="isEdit">
          <tiny-form-item label="状态" prop="status">
            <tiny-radio-group v-model="formData.status">
              <tiny-radio :label="0">启用</tiny-radio>
              <tiny-radio :label="1">禁用</tiny-radio>
            </tiny-radio-group>
          </tiny-form-item>
        </tiny-col>
        <tiny-col :span="12">
          <tiny-form-item label="用户类型" prop="user_type">
            <tiny-radio-group v-model="formData.user_type">
              <tiny-radio :label="0">管理员</tiny-radio>
              <tiny-radio :label="1">老师</tiny-radio>
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
  import { createAdminUser, updateAdminUser, getPinyin } from '@/api/user';
  import type { AdminUserRequest, AdminUserListItem } from '@/api/user';

  const visible = ref<boolean>(false);
  const props = defineProps<{
    visible: boolean;
    userData?: AdminUserListItem;
    isEdit?: boolean;
  }>();

  const emit = defineEmits(['update:visible', 'success']);
  const propsVisible = toRef(props, 'visible');
  const propsUserData = toRef(props, 'userData');
  const propsIsEdit = toRef(props, 'isEdit');

  const formRef = ref();
  const isEdit = ref(false);
  const loading = ref(false);
  const formData = reactive({
    id: undefined as number | undefined,
    name: '',
    username: '',
    status: 0,
    user_type: 1,
  });

  const rules = {
    name: [
      {
        required: true,
        message: '请输入姓名',
        trigger: 'blur',
      },
      {
        min: 1,
        max: 100,
        message: '姓名长度在 1 到 100 个字符',
        trigger: 'blur',
      },
    ],
    username: [
      {
        required: true,
        message: '请输入用户名',
        trigger: 'blur',
      },
      {
        min: 1,
        max: 50,
        message: '用户名长度在 1 到 50 个字符',
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
    user_type: [
      {
        required: true,
        message: '请选择用户类型',
        trigger: 'change',
      },
    ],
  };

  // 重置表单
  const resetForm = () => {
    formData.id = undefined;
    formData.name = '';
    formData.username = '';
    formData.status = 0;
    formData.user_type = 1;
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

  // 处理姓名输入
  const handleNameInput = async (val: string) => {
      try {
        const res = await getPinyin(formData.name);
        if (res.code === 200) {
          formData.username = res.data.pinyin;
        }
      } catch (error) {
        Modal.message({
          message: '获取姓名拼音失败',
          status: 'error',
        });
      }
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
            res = await updateAdminUser(submitData.id!, submitData);
          } else {
            res = await createAdminUser(submitData);
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
    () => propsUserData.value,
    async (val) => {
      if (val.id) {
        try {
          loading.value = true;
          // 直接使用传入的数据，不需要额外请求
          Object.assign(formData, val);
        } catch (error) {
          Modal.message({
            message: '获取管理员详情失败',
            status: 'error',
          });
          handleClose();
        } finally {
          loading.value = false;
        }
      } else {
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

  watch(propsIsEdit, (newValue) => {
    isEdit.value = newValue;
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

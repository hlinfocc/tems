<template>
  <tiny-drawer
    v-model:visible="visible"
    :title="isEdit ? '编辑评教问题集' : '新增评教问题集'"
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
          <tiny-form-item label="问题集名称" prop="name">
            <tiny-input
              v-model="formData.name"
              placeholder="请输入问题集名称"
              :maxlength="255"
              show-word-limit
            ></tiny-input>
          </tiny-form-item>
        </tiny-col>
        <tiny-col :span="24">
          <tiny-form-item label="备注" prop="remark">
            <tiny-input
              v-model="formData.remark"
              placeholder="请输入备注信息"
              type="textarea"
              :rows="4"
              :maxlength="500"
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
      <tiny-button type="info" @click="handleSubmit">
        {{ isEdit ? '更新' : '创建' }}
      </tiny-button>
    </template>
  </tiny-drawer>
</template>

<script lang="ts" setup>
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
  import type { EvaluationQuestionSetForm, EvaluationQuestionSetInfo } from '@/api/evaluationQuestionSet';
  import { addEvaluationQuestionSet, updateEvaluationQuestionSet } from '@/api/evaluationQuestionSet';

  const visible = ref<boolean>(false);
  const props = defineProps<{
    visible: boolean;
    questionSetData?: EvaluationQuestionSetInfo;
  }>();

  const emit = defineEmits(['update:visible', 'success']);
  const propsVisible = toRef(props, 'visible');
  const propsQuestionSetData = toRef(props, 'questionSetData');

  const formRef = ref();
  const isEdit = ref(false);
  const formData = reactive<EvaluationQuestionSetForm>({
    id: undefined,
    name: '',
    remark: '',
    status: 0,
  });

  const rules = {
    name: [
      {
        required: true,
        message: '请输入问题集名称',
        trigger: 'blur',
      },
    ],
  };

  // 重置表单
  const resetForm = () => {
    formData.id = undefined;
    formData.name = '';
    formData.remark = '';
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
  const handleSubmit = () => {
    formRef.value?.validate(async (valid: boolean) => {
      if (valid) {
        try {
          const submitData = { ...toRaw(formData) };
          let res;

          if (isEdit.value) {
            res = await updateEvaluationQuestionSet(submitData);
          } else {
            res = await addEvaluationQuestionSet(submitData);
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
        }
      }
    });
  };

  // 监听数据变化
  watch(
    () => propsQuestionSetData.value,
    (val) => {
      if (val && val.id) {
        isEdit.value = true;
        Object.assign(formData, val);
      } else {
        isEdit.value = false;
        resetForm();
      }
    },
    { immediate: true, deep: true },
  );
  watch(visible, (newValue, oldValue) => {
    emit('update:visible', visible.value);
  });
  watch(propsVisible, (newValue, oldValue) => {
    visible.value = propsVisible.value;
    console.log('watch propsVisible.value', propsVisible.value);
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
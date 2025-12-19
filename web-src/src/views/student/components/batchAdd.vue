<template>
  <tiny-drawer
    v-model:visible="visible"
    title="批量添加学生"
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
      :validate-type="`text`"
      :hide-required-asterisk="false"
      label-position="top"
      :popper-options="{ bubbling: true }"
    >
      <tiny-form-item
        label="学生信息"
        prop="studentInfos"
        extra="每行输入姓名，学号，班级（逗号分隔）"
        required
      >
        <tiny-input
          type="textarea"
          rows="20"
          v-model="formData.studentInfos"
          placeholder="请输入姓名，学号，班级（每行一个）"
        />
      </tiny-form-item>
    </tiny-form>
    <template #footer>
      <tiny-button @click="handleClose">取消</tiny-button>
      <tiny-button type="info" @click="handleSubmit"> 确定 </tiny-button>
    </template>
  </tiny-drawer>
</template>

<script lang="ts" setup>
  import { ref, reactive, watch, onMounted, toRaw, toRef } from 'vue';
  import {
    Form as TinyForm,
    FormItem as TinyFormItem,
    Input as TinyInput,
    Button as TinyButton,
    Drawer as TinyDrawer,
    Select as TinySelect,
    Option as TinyOption,
    Modal,
  } from '@opentiny/vue';
  import { addBatchStudent } from '@/api/student';
  import { queryClassList } from '@/api/class';
  import type { StudentInfo } from '@/api/student';

  // Props
  const props = defineProps<{
    visible: boolean;
    studentData?: any;
  }>();
  const visible = ref<boolean>(false);
  // Emits
  const emit = defineEmits(['update:visible', 'close', 'success']);
  const propsVisible = toRef(props, 'visible');
  // 表单引用
  const formRef = ref();

  // 班级列表
  const classOptions = ref<any[]>([]);

  // 表单数据
  const formData = reactive({
    studentID: '',
    studentName: '',
    phone: '',
    classID: undefined,
    className: '',
    password: '',
    studentInfos: '',
  });

  // 表单验证规则
  const rules = {
    studentInfos: [
      { required: true, message: '请输入学生信息', trigger: 'blur' },
    ],
    studentName: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
    phone: [
      {
        pattern: /^1[3-9]\d{9}$/,
        message: '请输入正确的手机号格式',
        trigger: 'blur',
      },
    ],
    classID: [{ required: true, message: '请选择班级', trigger: 'change' }],
    className: [{ required: true, message: '请输入班级名称', trigger: 'blur' }],
    password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  };

  // 重置表单
  function resetForm() {
    Object.assign(formData, {
      id: '',
      studentID: '',
      studentName: '',
      phone: '',
      classID: undefined,
      className: '',
      password: '',
      studentInfos: '',
    });
    if (formRef.value) {
      formRef.value.reset();
    }
  }

  // 关闭抽屉
  function handleClose() {
    console.log('handleClose', visible.value);
    emit('update:visible', false);
    emit('close');
    // 延时重置表单，确保抽屉完全关闭后再重置
    setTimeout(() => {
      resetForm();
    }, 300);
  }

  // 提交表单
  const handleSubmit = () => {
    if (!formRef.value) return;
    formRef.value?.validate(async (valid: boolean) => {
      // console.log('valid', valid);
      if (valid) {
        try {
          const submitData = { ...toRaw(formData) };
          console.log('submitData', submitData);
          if (!submitData.studentInfos) {
            return;
          }
          // 解析学生信息
          const studentInfos = submitData.studentInfos
            .split('\n')
            .map((line) => {
              const [name, id, className] = line.split(/[，,]/);
              return {
                studentName: name.trim(),
                studentID: id.trim(),
                className: className.trim(),
                classID: 0,
              };
            });
          // 检查是否有重复学号
          const hasDuplicateID = studentInfos.some(
            (item, index, arr) =>
              index !== arr.findIndex((t) => t.studentID === item.studentID),
          );
          if (hasDuplicateID) {
            Modal.message({
              message: '存在重复学号，请检查输入',
              status: 'error',
            });
            return;
          }

          // 新增模式
          const res = await addBatchStudent(studentInfos);
          if (res.code === 200) {
            Modal.message({
              message: res.msg || '添加成功',
              status: 'success',
            });
            emit('success');
            handleClose();
          }
        } catch (error) {
          console.error('表单验证失败:', error);
        }
      }
    });
  };

  // 加载班级列表
  async function loadClassList() {
    try {
      const res = await queryClassList();
      if (res.code === 200) {
        classOptions.value = res.data || [];
      }
    } catch (error) {
      console.error('加载班级列表失败:', error);
    }
  }
  watch(visible, (newValue) => {
    emit('update:visible', newValue);
  });

  watch(propsVisible, (newValue) => {
    visible.value = newValue;
  });
  // 生命周期
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

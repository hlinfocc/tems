<template>
  <tiny-drawer
    v-model:visible="visible"
    :title="studentData?.id ? '编辑学生' : '添加学生'"
    :width="`40%`"
    show-footer
    :mask-closable="false"
    @close="handleClose"
    @confirm="handleSubmit"
  >
    <tiny-form
      ref="form"
      :model="formData"
      :rules="rules"
      :validate-type="`text`"
      :hide-required-asterisk="false"
      label-position="top"
    >
      <tiny-form-item label="姓名" prop="studentName">
        <tiny-input v-model="formData.studentName" placeholder="请输入姓名" />
      </tiny-form-item>
      <tiny-form-item label="学号" prop="studentID">
        <tiny-input
          v-model="formData.studentID"
          :disabled="!!studentData?.id"
          placeholder="请输入学号"
        />
      </tiny-form-item>

      <tiny-form-item label="手机号" prop="phone">
        <tiny-input v-model="formData.phone" placeholder="请输入手机号" />
      </tiny-form-item>
      <tiny-form-item label="班级" prop="classID">
        <tiny-select v-model="formData.classID" placeholder="请选择班级">
          <tiny-option
            v-for="item in classOptions"
            :key="item.id"
            :label="item.className"
            :value="item.id"
            @click="formData.className = item.className"
          ></tiny-option>
        </tiny-select>
      </tiny-form-item>
    </tiny-form>
    <template #footer>
      <tiny-button @click="handleClose">取消</tiny-button>
      <tiny-button type="info" @click="handleSubmit">确定</tiny-button>
    </template>
  </tiny-drawer>
</template>

<script lang="ts" setup>
  import { ref, reactive, watch, onMounted, toRef } from 'vue';
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
  import { addStudent, updateStudent } from '@/api/student';
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
  const form = ref();

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
  });

  // 表单验证规则
  const rules = {
    studentID: [{ required: true, message: '请输入学号', trigger: 'blur' }],
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

  // 监听studentData变化，设置表单数据
  watch(
    () => props.studentData,
    (newVal) => {
      if (newVal && newVal.id) {
        // 编辑模式：填充表单数据
        Object.assign(formData, {
          id: newVal.id,
          studentID: newVal.studentID,
          studentName: newVal.studentName,
          phone: newVal.phone,
          classID: newVal.classID,
          className: newVal.className,
        });
      } else {
        // 新增模式：重置表单
        resetForm();
      }
    },
    { deep: true, immediate: true },
  );

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
    });
    if (form.value) {
      form.value.reset();
    }
  }

  // 关闭抽屉
  function handleClose() {
    emit('update:visible', false);
    emit('close');
    // 延时重置表单，确保抽屉完全关闭后再重置
    setTimeout(() => {
      resetForm();
    }, 300);
  }

  // 提交表单
  const handleSubmit = () => {
    if (!form.value) return;
    form.value?.validate(async (valid: boolean) => {
      console.log('valid', valid);
      if (valid) {
        try {
          const isEdit = !!props.studentData?.id;
          const submitData = { ...formData };

          // 如果是编辑模式，需要包含id
          if (isEdit && submitData.id) {
            const res = await updateStudent(submitData.id, submitData);
            if (res.code === 200) {
              Modal.message({
                message: res.msg || '更新成功',
                status: 'success',
              });
              emit('success');
              handleClose();
            }
          } else {
            // 新增模式
            const res = await addStudent(submitData);
            if (res.code === 200) {
              Modal.message({
                message: res.msg || '添加成功',
                status: 'success',
              });
              emit('success');
              handleClose();
            }
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
    loadClassList();
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

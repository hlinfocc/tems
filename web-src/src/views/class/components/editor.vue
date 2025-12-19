<template>
  <tiny-drawer
    v-model:visible="visible"
    :title="isEdit ? '编辑班级' : '新增班级'"
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
          <tiny-form-item label="班级名称" prop="className">
            <tiny-input
              v-model="formData.className"
              placeholder="请输入班级名称"
              :maxlength="50"
              show-word-limit
            ></tiny-input>
          </tiny-form-item>
        </tiny-col>
        <tiny-col :span="24">
          <tiny-form-item label="专业名称" prop="majorName">
            <tiny-input
              v-model="formData.majorName"
              placeholder="请输入专业名称"
              :maxlength="50"
              show-word-limit
            ></tiny-input>
          </tiny-form-item>
        </tiny-col>
        <tiny-col :span="24">
          <tiny-form-item label="年级" prop="grade">
            <tiny-input
              v-model="formData.grade"
              placeholder="请输入年级（如：2025）"
              :maxlength="20"
              show-word-limit
            ></tiny-input>
          </tiny-form-item>
        </tiny-col>
        <tiny-col :span="24">
          <tiny-form-item label="班主任" prop="headTeacherId">
            <tiny-select
              v-model="formData.headTeacherId"
              placeholder="请选择班主任"
              style="width: 100%"
              @change="handleTeacherChange"
            >
              <tiny-option
                v-for="user in userList"
                :key="user.id"
                :label="user.name"
                :value="user.id"
              ></tiny-option>
            </tiny-select>
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
    Select as TinySelect,
    Option as TinyOption,
    Modal,
  } from '@opentiny/vue';
  import type { ClassForm, ClassInfo } from '@/api/class';
  import { addClass, updateClass } from '@/api/class';
  import { getUserList, type UserListItem } from '@/api/user';

  const visible = ref<boolean>(false);
  const props = defineProps<{
    visible: boolean;
    classData?: ClassInfo;
  }>();

  const emit = defineEmits(['update:visible', 'success']);
  const propsVisible = toRef(props, 'visible');
  const propsClassData = toRef(props, 'classData');

  const formRef = ref();
  const isEdit = ref(false);
  const formData = reactive<ClassForm>({
    id: undefined,
    className: '',
    majorName: '',
    majorID: 0,
    grade: '',
    headTeacherId: 0,
    headTeacherName: '',
  });
  const userList = ref<UserListItem[]>([]);

  const rules = {
    className: [
      {
        required: true,
        message: '请输入班级名称',
        trigger: 'blur',
      },
    ],
    majorName: [
      {
        required: true,
        message: '请输入专业名称',
        trigger: 'blur',
      },
    ],
    grade: [
      {
        required: true,
        message: '请输入年级',
        trigger: 'blur',
      },
    ],
    headTeacherId: [
      {
        required: true,
        message: '请选择班主任',
        trigger: 'change',
      },
    ],
  };

  // 重置表单
  const resetForm = () => {
    formData.id = undefined;
    formData.className = '';
    formData.majorName = '';
    formData.majorID = 0;
    formData.grade = '';
    formData.headTeacherId = undefined;
    formData.headTeacherName = '';
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
            res = await updateClass(submitData);
          } else {
            res = await addClass(submitData);
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

  // 处理班主任选择变化
  const handleTeacherChange = (value: number) => {
    const selectedUser = userList.value.find((user) => user.id === value);
    if (selectedUser) {
      formData.headTeacherName = selectedUser.name;
    }
  };

  // 监听数据变化
  watch(
    () => propsClassData.value,
    (val) => {
      if (val && val.id) {
        isEdit.value = true;
        // 转换后端数据到表单格式
        formData.id = val.id;
        formData.className = val.className;
        formData.majorName = val.majorName;
        formData.majorID = val.majorID;
        formData.grade = val.grade;
        formData.headTeacherId = val.headTeacherId;
        formData.headTeacherName = val.headTeacherName;
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
  // watch(
  //   () => propsVisible.value,
  //   (val) => {
  //     console.log('watch propsVisible.value', val);
  //     visible.value = val;
  //   },
  // );
  // 获取用户列表
  const fetchUserList = async () => {
    try {
      const res = await getUserList({
        page: 0, 
        limit: 0, 
        isPage: 1,
      });
      if (res.code === 200) {
        userList.value = res.data.list;
      } else {
        Modal.message({
          message: '获取用户列表失败',
          status: 'error',
        });
      }
    } catch (error: any) {
      Modal.message({
        message: error.message || '获取用户列表失败',
        status: 'error',
      });
    }
  };

  onMounted(() => {
    visible.value = propsVisible.value;
    fetchUserList();
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

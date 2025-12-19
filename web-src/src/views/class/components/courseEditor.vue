<template>
  <tiny-drawer
    v-model:visible="visible"
    title="班级课程编辑"
    :width="`30%`"
    :mask-closable="true"
    @close="handleClose"
  >
    <div class="editor-container">
      <tiny-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-position="right"
        label-width="100px"
        size="medium"
      >
        <tiny-form-item label="课程" prop="courseId">
          <tiny-select
            v-model="formData.courseId"
            placeholder="请选择课程"
            @change="handleCourseChange"
          >
            <tiny-option
              v-for="course in courseList"
              :key="course.id"
              :label="course.courseName"
              :value="course.id"
            ></tiny-option>
          </tiny-select>
        </tiny-form-item>
        <tiny-form-item label="任课教师" prop="teacherId">
          <tiny-select
            v-model="formData.teacherId"
            placeholder="请选择任课教师"
            @change="handleTeacherChange"
          >
            <tiny-option
              v-for="teacher in teacherList"
              :key="teacher.id"
              :label="teacher.name"
              :value="teacher.id"
            ></tiny-option>
          </tiny-select>
        </tiny-form-item>
      </tiny-form>
      <div class="btn-group">
        <tiny-button @click="handleClose">取消</tiny-button>
        <tiny-button type="info" @click="handleSubmit">确定</tiny-button>
      </div>
    </div>
  </tiny-drawer>
</template>

<script lang="ts" setup>
  import { ref, reactive, toRefs, watch, onMounted } from 'vue';
  import {
    Drawer as TinyDrawer,
    Form as TinyForm,
    FormItem as TinyFormItem,
    Select as TinySelect,
    Option as TinyOption,
    Button as TinyButton,
    Modal,
  } from '@opentiny/vue';
  import type { ClassCourseForm } from '@/api/classCourse';
  import { addClassCourse, updateClassCourse } from '@/api/classCourse';
  import { queryCourseList } from '@/api/course';
  import { getUserList } from '@/api/user';

  const visible = ref<boolean>(false);
  const formRef = ref();
  const isEdit = ref(false);
  const courseList = ref<any[]>([]);
  const teacherList = ref<any[]>([]);

  const props = defineProps<{
    visible: boolean;
    courseData?: any;
    classId?: number;
    className?: string;
  }>();

  const emit = defineEmits(['update:visible', 'success']);

  const formData = reactive<ClassCourseForm>({
    id: undefined,
    classId: props.classId || 0,
    courseId: 0,
    teacherId: 0,
    className: '',
    courseName: '',
    teacherName: '',
  });

  const rules = {
    courseId: [
      {
        required: true,
        message: '请选择课程',
        trigger: 'change',
      },
    ],
    teacherId: [
      {
        required: true,
        message: '请选择任课教师',
        trigger: 'change',
      },
    ],
  };

  // 重置表单
  const resetForm = () => {
    formData.id = undefined;
    formData.classId = props.classId || 0;
    formData.courseId = undefined;
    formData.teacherId = undefined;
    formData.className = '';
    formData.courseName = '';
    formData.teacherName = '';
    if (formRef.value) {
      formRef.value.resetFields();
    }
  };

  // 关闭抽屉
  const handleClose = () => {
    visible.value = false;
    if (formRef.value) {
      formRef.value.resetFields();
    }
  };

  // 提交表单
  const handleSubmit = async () => {
    if (!formRef.value) return;
    try {
      await formRef.value.validate();
      let res;
      if (isEdit.value) {
        res = await updateClassCourse(formData.id!, formData);
      } else {
        res = await addClassCourse(formData);
      }
      if (res.code === 200) {
        Modal.message({
          message: res.msg,
          status: 'success',
        });
        visible.value = false;
        emit('success');
      } else {
        Modal.message({
          message: res.msg,
          status: 'error',
        });
      }
    } catch (error) {
      console.error('表单验证失败', error);
    }
  };

  // 加载课程列表
  const loadCourseList = async () => {
    try {
      const params = {
        page: 1,
        limit: 100,
      };
      const res: any = await queryCourseList(params);
      if (res.code === 200) {
        courseList.value = res.data;
      }
    } catch (error) {
      console.error('加载课程列表失败', error);
    }
  };

  // 加载教师列表
  const loadTeacherList = async () => {
    try {
      const params = {
        page: 1,
        limit: 100,
        isPage: 1,
      };
      const res: any = await getUserList(params);
      if (res.code === 200) {
        teacherList.value = res.data.list;
      }
    } catch (error) {
      console.error('加载教师列表失败', error);
    }
  };

  // 处理课程变化
  const handleCourseChange = (value: number) => {
    const course = courseList.value.find((item) => item.id === value);
    if (course) {
      formData.courseName = course.courseName;
    }
  };

  // 处理教师变化
  const handleTeacherChange = (value: number) => {
    const teacher = teacherList.value.find((item) => item.id === value);
    if (teacher) {
      formData.teacherName = teacher.name;
    }
  };

  // 监听课程数据变化
  watch(
    () => props.courseData,
    (val) => {
      if (val && (val.id || (val.classId && props.classId))) {
        isEdit.value = !!val.id;
        Object.assign(formData, val);
        // 如果有classId，优先使用props中的classId
        if (props.classId) {
          formData.classId = props.classId;
        }
        // 如果有className，优先使用props中的className
        if (props.className) {
          formData.className = props.className;
        }
      } else {
        isEdit.value = false;
        resetForm();
      }
    },
    { immediate: true, deep: true },
  );

  // 监听抽屉可见状态
  watch(
    () => props.visible,
    (newValue) => {
      visible.value = newValue;
      if (newValue) {
        // 当抽屉打开时，确保classId正确
        if (props.classId) {
          formData.classId = props.classId;
        }
      }
    },
  );

  // 监听classId变化
  watch(
    () => props.classId,
    (newValue) => {
      if (newValue) {
        formData.classId = newValue;
      }
    },
  );

  watch(visible, (newValue) => {
    emit('update:visible', newValue);
  });

  onMounted(async () => {
    visible.value = props.visible;
    // 加载课程和教师列表
    await Promise.all([loadCourseList(), loadTeacherList()]);
  });
</script>

<style scoped>
  .editor-container {
    padding: 20px;
  }

  .btn-group {
    display: flex;
    justify-content: flex-end;
    margin-top: 30px;
  }

  .btn-group :deep(.tiny-btn) {
    margin-left: 10px;
  }
</style>

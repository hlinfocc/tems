<template>
  <tiny-drawer
    v-model:visible="visible"
    :title="isEdit ? '编辑评教任务详情' : '新增评教任务详情'"
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
          <tiny-form-item label="班级" prop="classSelectVal" extra="选择班级将自动添加该班级下面的所有课程">
            <tiny-select
              v-model="formData.classSelectVal"
              placeholder="请选择班级"
              style="width: 100%"
              multiple
              show-alloption
              value-key="id"
              @change="handleClassSelect"
            >
              <tiny-option
                v-for="classItem in classList"
                :key="classItem.id"
                :label="classItem.className"
                :value="classItem"
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
  import { ref, reactive, watch, toRaw, onMounted } from 'vue';
  import {
    Drawer as TinyDrawer,
    Form as TinyForm,
    FormItem as TinyFormItem,
    Input as TinyInput,
    Button as TinyButton,
    Row as TinyRow,
    Col as TinyCol,
    Select as TinySelect,
    Option as TinyOption,
    Modal,
  } from '@opentiny/vue';
  import type { EvaluationTaskDetailForm } from '@/api/evaluationTaskDetail';
  import {
    batchAddEvaluationTaskDetail,
  } from '@/api/evaluationTaskDetail';
  import { queryClassList } from '@/api/class';
  import { queryCourseList } from '@/api/course';

  const visible = ref<boolean>(false);
  const formRef = ref();
  const isEdit = ref(false);
  const classList = ref<any[]>([]);
  const courseList = ref<any[]>([]);

  const props = defineProps<{
    visible: boolean;
    detailData?: any;
    taskId?: number;
  }>();

  const emit = defineEmits(['update:visible', 'success']);

  const formData = reactive<EvaluationTaskDetailForm>({
    id: undefined,
    taskId: props.taskId || 0,
    classId: 0,
    courseId: 0,
    className: '',
    courseName: '',
    userName: '',
    classSelectVal: [],
  });

  const rules = {
    classSelectVal: [
      {
        required: true,
        message: '请选择班级',
        trigger: 'change',
      },
    ],
    courseId: [
      {
        required: true,
        message: '请选择课程',
        trigger: 'change',
      },
    ],
    userName: [
      {
        required: true,
        message: '请输入任课教师名称',
        trigger: 'blur',
      },
    ],
  };

  // 重置表单
  const resetForm = () => {
    formData.id = undefined;
    formData.taskId = props.taskId || 0;
    formData.classId = 0;
    formData.courseId = 0;
    formData.className = '';
    formData.courseName = '';
    formData.userName = '';
    formData.classSelectVal = [];
    if (formRef.value) {
      formRef.value.resetFields();
    }
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
          const formDataRaw = { ...toRaw(formData) };
          let submitData = {
            taskId: formDataRaw.taskId,
            classes: [] as any[],
          };
          formDataRaw.classSelectVal.forEach((classItem: any) => {
            submitData.classes.push({
              classId: classItem.id,
              className: classItem.className,
            });
          });
          console.log(submitData);
          let res;

          res = await batchAddEvaluationTaskDetail(submitData);

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

  // 处理班级选择
  const handleClassSelect = (val: any, option: any) => {
    console.log(val, option);
    // formData.className = classItem.className;
  };

  // 处理课程选择
  const handleCourseSelect = (courseItem: any) => {
    formData.courseName = courseItem.courseName;
  };

  // 加载班级列表
  const loadClassList = async () => {
    try {
      const params = {
        page: 1,
        limit: 100,
      };
      const res: any = await queryClassList(params);
      if (res.code === 200) {
        classList.value = res.data;
      }
    } catch (error) {
      console.error('加载班级列表失败', error);
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

  // 监听详情数据变化
  watch(
    () => props.detailData,
    (val) => {
      if (val && (val.id || (val.taskId && props.taskId))) {
        isEdit.value = !!val.id;
        Object.assign(formData, val);
        // 如果有taskId，优先使用props中的taskId
        if (props.taskId) {
          formData.taskId = props.taskId;
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
        // 当抽屉打开时，确保taskId正确
        if (props.taskId) {
          formData.taskId = props.taskId;
        }
      }
    },
  );

  // 监听taskId变化
  watch(
    () => props.taskId,
    (newValue) => {
      if (newValue) {
        formData.taskId = newValue;
      }
    },
  );

  watch(visible, (newValue) => {
    emit('update:visible', newValue);
  });

  onMounted(async () => {
    visible.value = props.visible;
    // 加载班级和课程列表
    await Promise.all([loadClassList()]);
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

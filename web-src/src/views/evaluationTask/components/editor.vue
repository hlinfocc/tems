<template>
  <tiny-drawer
    v-model:visible="visible"
    :title="isEdit ? '编辑评教任务' : '新增评教任务'"
    :width="`50%`"
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
          <tiny-form-item label="任务名称" prop="taskName">
            <tiny-input
              v-model="formData.taskName"
              placeholder="请输入任务名称"
              :maxlength="255"
              show-word-limit
            ></tiny-input>
          </tiny-form-item>
        </tiny-col>
        <tiny-col :span="12">
          <tiny-form-item label="学年" prop="academicYear">
            <tiny-input
              v-model="formData.academicYear"
              placeholder="请输入学年（如：2023-2024）"
              :maxlength="50"
              show-word-limit
            ></tiny-input>
          </tiny-form-item>
        </tiny-col>
        <tiny-col :span="12">
          <tiny-form-item label="学期" prop="semester">
            <tiny-radio-group v-model="formData.semester">
              <tiny-radio :label="0">上学期</tiny-radio>
              <tiny-radio :label="1">下学期</tiny-radio>
            </tiny-radio-group>
          </tiny-form-item>
        </tiny-col>
        <tiny-col :span="12">
          <tiny-form-item label="问题集" prop="questionSetID">
            <tiny-select
              v-model="formData.questionSetID"
              placeholder="请选择问题集"
              filterable
              allow-create
              @change="handleQuestionSetIdChange"
            >
              <tiny-option
                v-for="item in questionSetListData"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              ></tiny-option>
            </tiny-select>
          </tiny-form-item>
        </tiny-col>
        <tiny-col :span="12">
          <tiny-form-item label="开始时间" prop="startTime">
            <tiny-date-picker
              v-model="formData.startTime"
              type="datetime"
              placeholder="请选择开始时间"
              style="width: 100%"
            ></tiny-date-picker>
          </tiny-form-item>
        </tiny-col>
        <tiny-col :span="12">
          <tiny-form-item label="结束时间" prop="endTime">
            <tiny-date-picker
              v-model="formData.endTime"
              type="datetime"
              placeholder="请选择结束时间"
              style="width: 100%"
            ></tiny-date-picker>
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
    DatePicker as TinyDatePicker,
    Modal,
  } from '@opentiny/vue';
  import type {
    EvaluationTaskForm,
    EvaluationTaskInfo,
  } from '@/api/evaluationTask';
  import {
    addEvaluationTask,
    updateEvaluationTask,
  } from '@/api/evaluationTask';
  import { queryEvaluationQuestionSetList } from '@/api/evaluationQuestionSet';

  const visible = ref<boolean>(false);
  const props = defineProps<{
    visible: boolean;
    taskData?: EvaluationTaskInfo;
  }>();

  const emit = defineEmits(['update:visible', 'success']);
  const propsVisible = toRef(props, 'visible');
  const propsTaskData = toRef(props, 'taskData');
  const questionSetListData = ref([]);

  const formRef = ref();
  const isEdit = ref(false);
  const formData = reactive<EvaluationTaskForm>({
    id: undefined,
    taskName: '',
    academicYear: '',
    semester: 0,
    questionSetID: undefined,
    questionSetName: '',
    startTime: '',
    endTime: '',
  });


  const rules = {
    taskName: [
      {
        required: true,
        message: '请输入任务名称',
        trigger: 'blur',
      },
    ],
    academicYear: [
      {
        required: true,
        message: '请输入学年',
        trigger: 'blur',
      },
    ],
    semester: [
      {
        required: true,
        message: '请选择学期',
        trigger: 'change',
      },
    ],
    questionSetID: [
      {
        required: true,
        message: '请选择问题集',
        trigger: 'change',
      },
    ],
    startTime: [
      {
        required: true,
        message: '请选择开始时间',
        trigger: 'change',
      },
    ],
    endTime: [
      {
        required: true,
        message: '请选择结束时间',
        trigger: 'change',
      },
      {
        validator: (rule: any, value: any, callback: any) => {
          if (value && formData.startTime) {
            if (new Date(value) <= new Date(formData.startTime)) {
              callback(new Error('结束时间必须晚于开始时间'));
            } else {
              callback();
            }
          } else {
            callback();
          }
        },
        trigger: 'change',
      },
    ],
  };

  async function fetchEQSetsData() {
    let params: any = {
      page: 1,
      limit: 10,
      isPage: 1,
      keyword: '',
      status: 0,
      currYear: 1,
    };
    const res: any = await queryEvaluationQuestionSetList(params);
    console.log('res:', res);
    if (res.code === 200) {
      questionSetListData.value = res.data;
    }
  }
  const handleQuestionSetIdChange = (val: any) => {
    console.log('val:', val);
    if (val>0) {
      formData.questionSetName = questionSetListData.value.find(
        (item: any) => item.id === val
      )?.name || '';
    } else {
      formData.questionSetName = '';
    }
    console.log('formData:', formData);
  };

  // 监听visible变化
  watch(propsVisible, (newVal) => {
    visible.value = newVal;
  });
  watch(propsTaskData, (newVal) => {
    console.log('============>>>>>>>>', newVal);
    if (newVal.id) {
      // 编辑模式，填充表单数据
      isEdit.value = true;
      fillFormData(propsTaskData.value);
    } else {
      // 新增模式，重置表单
      isEdit.value = false;
      resetForm();
    }
  });

  // 填充表单数据
  const fillFormData = (taskData: EvaluationTaskInfo) => {
    formData.id = taskData.id;
    formData.taskName = taskData.taskName;
    formData.academicYear = taskData.academicYear;
    formData.semester = taskData.semester;
    formData.questionSetID = taskData.questionSetID;
    formData.questionSetName = taskData.questionSetName;
    formData.startTime = taskData.startTime;
    formData.endTime = taskData.endTime;
  };

  // 重置表单
  const resetForm = () => {
    formData.id = undefined;
    formData.taskName = '';
    formData.academicYear = '';
    formData.semester = 0;
    formData.questionSetID = undefined;
    formData.questionSetName = '';
    formData.startTime = '';
    formData.endTime = '';
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
          // 创建提交数据，确保字段名正确
          const submitData = {
            ...toRaw(formData),
            // 确保字段名与后端匹配
            taskName: formData.taskName,
            academicYear: formData.academicYear,
            semester: formData.semester,
            questionSetID: formData.questionSetID,
            startTime: formData.startTime,
            endTime: formData.endTime,
          };
          let res;

          if (isEdit.value) {
            res = await updateEvaluationTask(submitData);
          } else {
            res = await addEvaluationTask(submitData);
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
        } catch (error) {
          Modal.message({
            message: '操作失败，请稍后重试',
            status: 'error',
          });
        }
      }
    });
  };

  // 获取问题集列表（实际项目中应该从API获取）
  const fetchQuestionSets = async () => {
    // 模拟API调用
    // 实际项目中应该调用后端接口获取问题集列表
    // 这里使用模拟数据
    console.log('获取问题集列表');
  };

  onMounted(async () => {
    await fetchEQSetsData();
    fetchQuestionSets();
  });
</script>

<style scoped lang="less">
  .tiny-form-item {
    margin-bottom: 20px;
  }
  .tiny-select {
    width: 100%;
  }
</style>

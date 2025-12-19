<template>
  <tiny-drawer
    v-model:visible="visible"
    :title="isEdit ? '编辑评教问题详情' : '新增评教问题详情'"
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
          <tiny-form-item label="问题集" prop="questionSetId">
            <tiny-select
              v-model="formData.questionSetId"
              placeholder="请选择问题集"
            >
              <tiny-option
                v-for="value in questionSetListData"
                :key="value.id"
                :label="value.name"
                :value="value.id"
              ></tiny-option>
            </tiny-select>
          </tiny-form-item>
        </tiny-col>
        <tiny-col :span="24">
          <tiny-form-item label="问题标题" prop="title">
            <tiny-input
              v-model="formData.title"
              placeholder="请输入问题标题"
              :maxlength="500"
              show-word-limit
            ></tiny-input>
          </tiny-form-item>
        </tiny-col>
        <tiny-col :span="24">
          <tiny-form-item label="问题类型" prop="questionType">
            <tiny-radio-group v-model="formData.questionType">
              <tiny-radio :label="1">单选</tiny-radio>
              <tiny-radio :label="2">多选</tiny-radio>
            </tiny-radio-group>
          </tiny-form-item>
        </tiny-col>
        <tiny-col :span="24">
          <tiny-form-item label="选项">
            <div
              v-for="(optionItem, index) in optionsList"
              :key="index"
              class="option-item"
            >
              <tiny-input
                v-model="optionItem.value"
                placeholder="请输入选项内容"
                :maxlength="200"
                show-word-limit
                style="margin-bottom: 10px"
              >
                <template #prepend>{{ optionItem.key }}.</template>
              </tiny-input>
              <tiny-button
                type="text"
                danger
                @click="removeOption(index)"
                v-if="optionsList.length > 1"
                >删除</tiny-button
              >
            </div>
            <tiny-button
              type="info"
              :size="`mini`"
              @click="addOption"
              class="add-option-btn btn-mini"
            >
              + 添加选项
            </tiny-button>
            <tiny-button
              type="info"
              size="mini"
              @click="batchAddOptionClick"
              class="add-option-btn btn-mini"
              plain
            >
              + 批量添加选项
            </tiny-button>
          </tiny-form-item>
        </tiny-col>
        <tiny-col :span="24">
          <tiny-form-item label="备注" prop="remark">
            <tiny-input
              v-model="formData.remark"
              placeholder="请输入备注"
              type="textarea"
              :rows="3"
              :maxlength="500"
              show-word-limit
            ></tiny-input>
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
    <tiny-dialog-box
      v-model:visible="boxVisibility"
      title="批量添加选项"
      width="30%"
      :modal="false"
    >
      <div>
        <tiny-input
          v-model="batchAddOperation"
          placeholder="请输入批量添加选项的操作，每个选项占一行，只需要输入内容，不要输入选项的序号"
          type="textarea"
          :rows="5"
        ></tiny-input>
      </div>
      <template #footer>
        <tiny-button @click="boxVisibility = false" round>取 消</tiny-button>
        <tiny-button type="info" @click="batchAddOptionsHandle" round
          >确 定</tiny-button
        >
      </template>
    </tiny-dialog-box>
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
    TinyDialogBox,
  } from '@opentiny/vue';
  import type {
    QuestionDetailForm,
    QuestionDetailInfo,
  } from '@/api/evaluation-question-detail';
  import {
    addQuestionDetail,
    updateQuestionDetail,
  } from '@/api/evaluation-question-detail';
  import { queryEvaluationQuestionSetList } from '@/api/evaluationQuestionSet';

  const visible = ref<boolean>(false);
  const boxVisibility = ref(false);
  const batchAddOperation = ref('');
  const props = defineProps<{
    visible: boolean;
    questionData?: QuestionDetailInfo;
  }>();

  const emit = defineEmits(['update:visible', 'success']);
  const propsVisible = toRef(props, 'visible');
  const propsQuestionData = toRef(props, 'questionData');

  // 问题集列表数据
  const questionSetListData = ref([]);

  const formRef = ref();
  const isEdit = ref(false);
  const formData = reactive<QuestionDetailForm>({
    id: undefined,
    questionSetId: undefined,
    title: '',
    questionType: 0,
    options: '',
    remark: '',
  });

  // 用于存储选项列表
  const optionsList = ref<string[]>([]);

  const rules = {
    questionSetId: [
      {
        required: true,
        message: '请输入问题集ID',
        trigger: 'blur',
      },
    ],
    title: [
      {
        required: true,
        message: '请输入问题标题',
        trigger: 'blur',
      },
    ],
    questionType: [
      {
        required: true,
        message: '请选择问题类型',
        trigger: 'change',
      },
    ],
  };

  // 添加选项
  const addOption = () => {
    // 生成新的选项键（例如：A、B、C...）
    const newKey = String.fromCharCode(65 + optionsList.value.length);
    optionsList.value.push({ key: newKey, value: '', ans: 0 });
  };

  // 批量添加选项
  const batchAddOptionClick = () => {
    batchAddOperation.value = '';
    boxVisibility.value = true;
  };
  const batchAddOptionsHandle = () => {
    if (!batchAddOperation.value) {
      Modal.message({
        message: '请输入批量添加选项的内容',
        status: 'warning',
      });
      return;
    }

    // 按行分割操作
    const lines = batchAddOperation.value.split('\n');
    lines.forEach((line) => {
      line = line.trim();
      if (line) {
        // 生成新的选项键（例如：A、B、C...）
        const newKey = String.fromCharCode(65 + optionsList.value.length);
        optionsList.value.push({ key: newKey, value: line });
      }
    });

    // 清空输入框
    batchAddOperation.value = '';
    boxVisibility.value = false;
  };

  // 删除选项
  const removeOption = (index: number) => {
    optionsList.value.splice(index, 1);
    // 更新选项键
    optionsList.value.forEach((opt, i) => {
      opt.key = String.fromCharCode(65 + i);
    });
  };

  // 将选项列表转换为JSON字符串
  const convertOptionsToJson = () => {
    // 过滤掉空选项
    const validOptions = optionsList.value.filter((opt) => opt.value.trim() !== '');
    return JSON.stringify(validOptions);
  };

  // 将JSON字符串转换为选项列表
  const convertJsonToOptions = (jsonStr: string) => {
    if (!jsonStr) {
      optionsList.value = [{ key: 'A', value: '', ans: 0 }];
      return;
    }
    try {
      const options = JSON.parse(jsonStr);
      optionsList.value = Array.isArray(options)
        ? options
        : [{ key: 'A', value: '', ans: 0 }];
    } catch {
      optionsList.value = [{ key: 'A', value: '', ans: 0 }];
    }
  };

  // 重置表单
  const resetForm = () => {
    formData.id = undefined;
    formData.questionSetId = undefined;
    formData.title = '';
    formData.questionType = 1;
    formData.options = '';
    formData.remark = '';
    optionsList.value = [{ key: 'A', value: '', ans: 0 }];
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
    // 验证选项
    const validOptions = optionsList.value.filter((opt) => opt.value.trim() !== '');
    if (validOptions.length === 0) {
      Modal.message({
        message: '请至少添加一个有效选项',
        status: 'warning',
      });
      return;
    }

    formRef.value?.validate(async (valid: boolean) => {
      if (valid) {
        try {
          // 转换选项为JSON字符串
          formData.options = convertOptionsToJson();
          const submitData = { ...toRaw(formData) };
          let res;

          if (isEdit.value) {
            res = await updateQuestionDetail(submitData);
          } else {
            res = await addQuestionDetail(submitData);
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
  async function fetchEQSetsData() {
    let params: any = {
      page: 1,
      limit: 10,
      isPage: 1,
      keyword: '',
      status: 0,
      currYear:1,
    };
    const res: any = await queryEvaluationQuestionSetList(params);
    console.log('res:', res);
    if (res.code === 200) {
      questionSetListData.value = res.data;
    }
  }
  // 监听数据变化
  watch(
    () => propsQuestionData.value,
    (val) => {
      if (val && val.id) {
        isEdit.value = true;
        Object.assign(formData, val);
        // 转换JSON字符串为选项列表
        convertJsonToOptions(val.options || '');
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
    fetchEQSetsData();
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

  .option-item {
    display: flex;
    align-items: center;
  }

  .option-item :deep(.tiny-input) {
    flex: 1;
    margin-right: 10px;
  }

  .add-option-btn {
    margin-top: 10px;
  }
  .btn-mini {
    font-size: var(--tv-Button-font-size-mini);
    height: var(--tv-Button-height-mini);
    padding: 0 var(--tv-Button-padding-x-mini);
    min-width: var(--tv-Button-min-width-mini);
  }
</style>

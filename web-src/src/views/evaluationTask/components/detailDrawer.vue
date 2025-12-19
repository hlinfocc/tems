<template>
  <tiny-drawer
    v-model:visible="visible"
    title="评教任务详情"
    :width="`85%`"
    :mask-closable="false"
    @close="handleClose"
  >
    <div class="detail-container">
      <div class="task-info">
        <h3>{{ taskData.taskName }} - 详情列表</h3>
        <div class="task-meta">
          <span>学年：{{ taskData.academicYear }}</span>
          <span>学期：{{ taskData.semester === 0 ? '上学期' : '下学期' }}</span>
          <span>问题集：{{ taskData.questionSetName }}</span>
        </div>
      </div>

      <div class="detail-actions">
        <tiny-button type="info" @click="handleAddDetail">添加</tiny-button>
        <!-- <tiny-button
          type="danger"
          @click="handleBatchDelete"
          :disabled="selectedRows.length === 0"
          >批量删除</tiny-button
        > -->
      </div>

      <tiny-grid
        ref="detailGrid"
        :fetch-data="fetchDataOption"
        :pager="pagerConfig"
        :loading="loading"
        seq-serial
        size="medium"
        :auto-resize="true"
        row-id="id"
        :show-checkbox="true"
        @select="handleSelect"
        @select-all="handleSelectAll"
      >
        <tiny-grid-column type="index" width="60"></tiny-grid-column>
        <tiny-grid-column field="className" title="班级名称" align="center">
        </tiny-grid-column>
        <tiny-grid-column field="courseName" title="课程名称" align="center">
        </tiny-grid-column>
        <tiny-grid-column field="userName" title="任课教师" align="center">
        </tiny-grid-column>
        <tiny-grid-column field="updatedAt" title="更新时间" align="center">
        </tiny-grid-column>
        <tiny-grid-column title="操作" align="center">
          <template v-slot="data">
            <a class="operation-item" @click="handleDeleteDetail(data.row.id)">
              删除
            </a>
          </template>
        </tiny-grid-column>
      </tiny-grid>
    </div>

    <!-- 详情编辑子抽屉 -->
    <detailEditor
      v-model:visible="editorVisible"
      :detail-data="selectedDetail"
      :task-id="taskData.id"
      @success="handleDetailSuccess"
    />
  </tiny-drawer>
</template>

<script lang="ts" setup>
  import { ref, reactive, toRefs, watch, onMounted } from 'vue';
  import {
    Drawer as TinyDrawer,
    Grid as TinyGrid,
    GridColumn as TinyGridColumn,
    Button as TinyButton,
    Pager as TinyPager,
    Modal,
  } from '@opentiny/vue';
  import type { EvaluationTaskInfo } from '@/api/evaluationTask';
  import type {
    EvaluationTaskDetailInfo,
    QueryParams,
  } from '@/api/evaluationTaskDetail';
  import {
    getTaskDetailsByTaskId,
    deleteEvaluationTaskDetail,
  } from '@/api/evaluationTaskDetail';
  import DetailEditor from './detailEditor.vue';

  const visible = ref<boolean>(false);
  const editorVisible = ref<boolean>(false);
  const selectedDetail = ref<any>({});
  const selectedRows = ref<any[]>([]);

  const props = defineProps<{
    visible: boolean;
    taskData: EvaluationTaskInfo;
  }>();

  const emit = defineEmits(['update:visible', 'success']);
  const propsVisible = toRefs(props).visible;
  const propsTaskData = toRefs(props).taskData;

  // 加载效果
  const state = reactive<{
    loading: boolean;
  }>({
    loading: false,
  });

  const pagerConfig = reactive({
    component: TinyPager,
    attrs: {
      currentPage: 1,
      pageSize: 10,
      pageSizes: [10, 20],
      total: 10,
      layout: 'total, prev, pager, next, jumper, sizes',
    },
  });

  const { loading } = toRefs(state);
  const detailGrid = ref();
  const taskData = ref<EvaluationTaskInfo>({} as EvaluationTaskInfo);

  // 请求数据接口方法
  async function fetchData(
    params: QueryParams = {
      page: 1,
      limit: 10,
      taskId: undefined,
    },
  ) {
    if (!taskData.value.id) return { result: [], page: { total: 0 } };

    const queryParams = {
      ...params,
      taskId: taskData.value.id,
    };

    state.loading = true;
    try {
      const res: any = await getTaskDetailsByTaskId(queryParams);
      if (res.code !== 200) {
        Modal.message({
          message: res.msg,
          status: 'error',
        });
        return {
          result: [],
          page: { total: 0 },
        };
      }
      return {
        result: res.data,
        page: { total: res.count },
      };
    } finally {
      state.loading = false;
    }
  }

  const fetchDataOption = reactive({
    api: ({ page }: any) => {
      const { currentPage, pageSize } = page;

      return fetchData({
        page: currentPage,
        limit: pageSize,
      });
    },
  });

  // 处理选中行
  const handleSelect = (rows: any[]) => {
    selectedRows.value = rows;
  };

  // 处理全选
  const handleSelectAll = (rows: any[]) => {
    selectedRows.value = rows;
  };

  // 添加详情
  const handleAddDetail = () => {
    selectedDetail.value = { taskId: taskData.value.id };
    editorVisible.value = true;
  };

  // 编辑详情
  const handleEditDetail = (row: EvaluationTaskDetailInfo) => {
    selectedDetail.value = { ...row };
    editorVisible.value = true;
  };

  // 删除详情
  const handleDeleteDetail = (id: number) => {
    Modal.confirm('您确定要删除该评教任务详情吗？').then((rs: any) => {
      deleteEvaluationTaskDetail(id).then((res: any) => {
        if (res.code === 200) {
          reloadGrid();
          Modal.message({
            message: res.msg,
            status: 'success',
          });
        } else {
          Modal.message({
            message: res.msg,
            status: 'error',
          });
        }
      });
    });
  };

  // 批量删除
  const handleBatchDelete = () => {
    if (selectedRows.value.length === 0) {
      Modal.message({
        message: '请选择要删除的详情',
        status: 'warning',
      });
      return;
    }

    Modal.confirm(
      `您确定要删除选中的${selectedRows.value.length}条评教任务详情吗？`,
    ).then(async (rs: any) => {
      try {
        // 使用Promise.all并行执行删除操作
        await Promise.all(
          selectedRows.value.map((row) => deleteEvaluationTaskDetail(row.id)),
        );
        reloadGrid();
        selectedRows.value = [];
        Modal.message({
          message: '删除成功',
          status: 'success',
        });
      } catch (error) {
        Modal.message({
          message: '删除失败',
          status: 'error',
        });
      }
    });
  };

  // 详情操作成功回调
  const handleDetailSuccess = () => {
    reloadGrid();
    emit('success');
  };

  // 刷新表格
  const reloadGrid = () => {
    detailGrid?.value.handleFetch('reload');
    fetchData();
  };

  // 关闭抽屉
  const handleClose = () => {
    visible.value = false;
    emit('update:visible', false);
    selectedRows.value = [];
  };

  // 监听任务数据变化
  watch(
    () => propsTaskData.value,
    (val) => {
      if (val && val.id) {
        taskData.value = { ...val };
        if (visible.value) {
          reloadGrid();
        }
      }
    },
    { immediate: true, deep: true },
  );

  // 监听抽屉可见状态
  watch(propsVisible, (newValue) => {
    visible.value = newValue;
    if (newValue && taskData.value.id) {
      reloadGrid();
    }
  });

  watch(visible, (newValue) => {
    emit('update:visible', newValue);
  });

  onMounted(() => {
    visible.value = propsVisible.value;
  });
</script>

<style scoped>
  .detail-container {
    height: calc(100vh - 300px);
    display: flex;
    flex-direction: column;
  }

  .task-info {
    margin-bottom: 20px;
  }

  .task-info h3 {
    margin-bottom: 10px;
    color: #333;
  }

  .task-meta {
    display: flex;
    gap: 20px;
    color: #666;
    font-size: 14px;
  }

  .detail-actions {
    display: flex;
    gap: 10px;
    margin-bottom: 20px;
  }

  :deep(.tiny-grid) {
    flex: 1;
    overflow: hidden;
  }
</style>

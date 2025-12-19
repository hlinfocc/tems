<template>
  <div class="container-list">
    <Breadcrumb :items="['评教管理', '评教任务列表']" :is-i18n="false" />
    <div class="contain">
      <tiny-form
        :model="filterOptions"
        label-position="right"
        label-width="100px"
        class="filter-form"
        size="small"
      >
        <tiny-row :flex="true" justify="center" class="col">
          <tiny-col :span="4" label-width="100px">
            <tiny-form-item label="任务名称">
              <tiny-input
                v-model="filterOptions.keyword"
                placeholder="请输入任务名称"
              ></tiny-input>
            </tiny-form-item>
          </tiny-col>
          <tiny-col :span="4" label-width="100px">
            <tiny-form-item label="学年">
              <tiny-input
                v-model="filterOptions.academicYear"
                placeholder="请输入学年"
              ></tiny-input>
            </tiny-form-item>
          </tiny-col>
          <tiny-col :span="4" label-width="100px">
            <tiny-form-item label="学期">
              <tiny-select
                v-model="filterOptions.semester"
                placeholder="请选择学期"
              >
                <tiny-option :label="'上学期'" :value="0"></tiny-option>
                <tiny-option :label="'下学期'" :value="1"></tiny-option>
              </tiny-select>
            </tiny-form-item>
          </tiny-col>
          <tiny-col :span="4">
            <tiny-form-item>
              <div class="search-btn">
                <tiny-button type="info" @click="reloadGrid">
                  搜索
                </tiny-button>
                <tiny-button @click="handleFormReset"> 重置 </tiny-button>
              </div>
            </tiny-form-item>
          </tiny-col>
        </tiny-row>
      </tiny-form>
      <div class="bottom-line">
        <hr />
      </div>
      <tiny-fullscreen
        :teleport="true"
        :page-only="true"
        :z-index="999"
        :fullscreen="fullscreen"
        @update:fullscreen="fullscreen = $event"
      >
        <div class="tiny-fullscreen-scroll">
          <div class="tiny-fullscreen-wrapper">
            <tiny-grid
              ref="taskGrid"
              :fetch-data="fetchDataOption"
              :pager="pagerConfig"
              :loading="loading"
              seq-serial
              size="medium"
              :auto-resize="true"
              row-id="id"
            >
              <template #toolbar>
                <tiny-grid-toolbar>
                  <template #buttons>
                    <div class="btn">
                      <tiny-button
                        type="info"
                        @click="
                          () => {
                            selectedTask = {};
                            editorVisible = true;
                          }
                        "
                      >
                        添加
                      </tiny-button>
                      <div class="screen">
                        <img
                          v-if="!fullscreen"
                          src="@/assets/images/screen-out.png"
                          class="screen-image"
                          @click="toggle"
                        />
                        <img
                          v-if="fullscreen"
                          src="@/assets/images/screen-in.png"
                          class="screen-image"
                          @click="toggle"
                        />
                        <span @click="toggle">
                          {{ fullscreen ? '退出全屏' : '全屏' }}
                        </span>
                      </div>
                    </div>
                  </template>
                </tiny-grid-toolbar>
              </template>
              <tiny-grid-column type="index" width="60"></tiny-grid-column>
              <tiny-grid-column
                field="taskName"
                title="任务名称"
                align="center"
              >
              </tiny-grid-column>
              <tiny-grid-column
                field="academicYear"
                title="学年"
                align="center"
              >
              </tiny-grid-column>
              <tiny-grid-column field="semester" title="学期" align="center">
                <template #default="{ row }">
                  <tiny-tag :type="row.semester === 0 ? 'success' : 'primary'">
                    {{ row.semester === 0 ? '上学期' : '下学期' }}
                  </tiny-tag>
                </template>
              </tiny-grid-column>
              <tiny-grid-column
                field="questionSetName"
                title="问题集名称"
                align="center"
              >
              </tiny-grid-column>
              <tiny-grid-column
                field="startTime"
                title="开始时间"
                align="center"
              >
              </tiny-grid-column>
              <tiny-grid-column field="endTime" title="结束时间" align="center">
              </tiny-grid-column>
              <tiny-grid-column title="状态" align="center">
                <template #default="{ row }">
                  <tiny-tag :type="getStatusType(row)">
                    {{ getStatusText(row) }}
                  </tiny-tag>
                </template>
              </tiny-grid-column>
              <tiny-grid-column
                field="updatedAt"
                title="更新时间"
                align="center"
              >
              </tiny-grid-column>
              <tiny-grid-column title="操作" align="center">
                <template v-slot="data">
                  <a class="operation-item" @click="handleDetail(data.row)">
                    详情
                  </a>
                  <a class="operation-item" @click="handleEdit(data.row)">
                    编辑
                  </a>
                  <a class="operation-item" @click="handleDelete(data.row.id)">
                    删除
                  </a>
                </template>
              </tiny-grid-column>
            </tiny-grid>
          </div>
        </div>
      </tiny-fullscreen>
    </div>
    <editor
      v-model:visible="editorVisible"
      :task-data="selectedTask"
      @success="addCallback"
    />
    <detailDrawer
      v-model:visible="detailVisible"
      :task-data="selectedTask"
      @success="addCallback"
    />
  </div>
</template>

<script lang="ts" setup>
  import type { QueryParams } from '@/api/evaluationTask';
  import { ref, reactive, toRefs, onMounted, toRaw } from 'vue';
  import {
    Grid as TinyGrid,
    GridColumn as TinyGridColumn,
    GridToolbar as TinyGridToolbar,
    Form as TinyForm,
    FormItem as TinyFormItem,
    Input as TinyInput,
    Button as TinyButton,
    Row as TinyRow,
    Col as TinyCol,
    Select as TinySelect,
    Option as TinyOption,
    Pager as TinyPager,
    Fullscreen as TinyFullscreen,
    Modal,
    TinyTag,
  } from '@opentiny/vue';
  import {
    queryEvaluationTaskList,
    deleteEvaluationTask,
  } from '@/api/evaluationTask';
  import Breadcrumb from '@/components/breadcrumb/index.vue';
  import Editor from './components/editor.vue';
  import DetailDrawer from './components/detailDrawer.vue';

  const editorVisible = ref(false);
  const detailVisible = ref(false);
  const selectedTask = ref<any>({});

  // 初始化请求数据
  interface FilterOptions {
    keyword: string;
    academicYear: string;
    semester: number | null;
  }

  // 加载效果
  const state = reactive<{
    loading: boolean;
    filterOptions: FilterOptions;
  }>({
    loading: false,
    filterOptions: {
      keyword: '',
      academicYear: '',
      semester: null,
    },
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

  let tableData = ref([]);
  const taskGrid = ref();
  const { loading, filterOptions } = toRefs(state);

  // 请求数据接口方法
  async function fetchData(
    params: QueryParams = {
      page: 1,
      limit: 10,
      keyword: '',
      academic_year: '',
      semester: undefined,
    },
  ) {
    console.log('filterOptions.value:', filterOptions);

    const { ...rest } = filterOptions;
    // 移除null值，避免发送到后端
    const cleanRest = { ...rest };
    if (cleanRest.semester === null) {
      delete cleanRest.semester;
    }

    const queryParams = {
      ...params,
      ...cleanRest,
    };

    state.loading = true;
    try {
      const res: any = await queryEvaluationTaskList(queryParams);
      console.log('res:', res);
      tableData.value = res.data;
      return {
        result: res.data,
        page: { total: res.count },
      };
    } finally {
      state.loading = false;
    }
  }

  // 获取任务状态类型
  const getStatusType = (row: any) => {
    const now = new Date();
    const startTime = new Date(row.startTime);
    const endTime = new Date(row.endTime);

    if (now < startTime) {
      return 'warning'; // 未开始
    }
    if (now >= startTime && now <= endTime) {
      return 'success'; // 进行中
    }
    return 'danger'; // 已结束
  };

  // 获取任务状态文本
  const getStatusText = (row: any) => {
    const now = new Date();
    const startTime = new Date(row.startTime);
    const endTime = new Date(row.endTime);

    if (now < startTime) {
      return '未开始';
    }
    if (now >= startTime && now <= endTime) {
      return '进行中';
    }
    return '已结束';
  };

  const addCallback = () => {
    reloadGrid();
  };

  const handleDetail = (e: any) => {
    detailVisible.value = true;
    selectedTask.value = e;
  };

  const fetchDataOption = reactive({
    api: ({ page }: any) => {
      const { currentPage, pageSize } = page;

      return fetchData({
        page: currentPage,
        limit: pageSize,
      });
    },
  });

  const handleDelete = (id: number) => {
    Modal.confirm('您确定要删除该评教任务吗？').then((rs: any) => {
      deleteEvaluationTask(id).then((res: any) => {
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

  const handleEdit = (e: any) => {
    editorVisible.value = true;
    console.log('task data:', toRaw(e));
    selectedTask.value = e;
  };

  // form的button
  function reloadGrid() {
    taskGrid?.value.handleFetch('reload');
  }

  function handleFormReset() {
    state.filterOptions = {
      keyword: '',
      academicYear: '',
      semester: null,
    };
    reloadGrid();
  }

  // 全屏缩放设置
  const fullscreen = ref(false);
  const toggle = () => {
    fullscreen.value = !fullscreen.value;
  };

  onMounted(() => {
    reloadGrid();
  });
</script>

<style scoped lang="less">
  @import '@/assets/style/page-table.less';
</style>

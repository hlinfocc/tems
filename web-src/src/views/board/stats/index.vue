<template>
  <div class="container-list">
    <Breadcrumb :items="['统计分析', '评教详情统计']" :is-i18n="false" />
    <div class="contain">
      <tiny-form
        :model="filterOptions"
        label-position="right"
        label-width="100px"
        class="filter-form"
        size="small"
      >
        <tiny-row :flex="true" justify="center" class="col">
          <tiny-col :span="5" label-width="100px">
            <tiny-form-item label="任务">
              <tiny-select
                v-model="filterOptions.taskId"
                filterable
                remote
                :remote-method="loadTaskData"
                :loading="taskLoading"
                loading-text="Loading..."
                placeholder="请选择任务"
              >
                <tiny-option
                  v-for="item in taskDataList"
                  :key="item.id"
                  :label="item.taskName"
                  :value="item.id"
                >
                </tiny-option>
              </tiny-select>
            </tiny-form-item>
          </tiny-col>
          <tiny-col :span="5" label-width="100px">
            <tiny-form-item label="班级">
              <tiny-select
                v-model="filterOptions.classId"
                filterable
                remote
                :remote-method="loadClassesData"
                :loading="classesLoading"
                loading-text="Loading..."
                placeholder="请选择班级"
              >
                <tiny-option
                  v-for="item in classesDataList"
                  :key="item.id"
                  :label="item.className"
                  :value="item.id"
                >
                </tiny-option>
              </tiny-select>
            </tiny-form-item>
          </tiny-col>
          <tiny-col :span="5" label-width="100px">
            <tiny-form-item label="课程">
              <tiny-select
                v-model="filterOptions.courseId"
                filterable
                remote
                :remote-method="loadCourseData"
                :loading="courseLoading"
                loading-text="Loading..."
                placeholder="请选择课程"
              >
                <tiny-option
                  v-for="item in courseDataList"
                  :key="item.id"
                  :label="item.courseName"
                  :value="item.id"
                >
                </tiny-option>
              </tiny-select>
            </tiny-form-item>
          </tiny-col>
          <tiny-col :span="5" label-width="100px">
            <tiny-form-item label="教师">
              <tiny-select
                v-model="filterOptions.userId"
                filterable
                remote
                :remote-method="loadAdminUserData"
                :loading="adminUserLoading"
                loading-text="Loading..."
                placeholder="请选择教师"
              >
                <tiny-option
                  v-for="item in adminUserDataList"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                >
                </tiny-option>
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
                      <tiny-button type="info" @click="exportStats">
                        导出
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
                field="task_name"
                title="任务名称"
                align="center"
              ></tiny-grid-column>
              <tiny-grid-column
                field="academic_year"
                title="学年"
                align="center"
              ></tiny-grid-column>
              <tiny-grid-column
                field="semester"
                title="学期"
                align="center"
              ></tiny-grid-column>
              <tiny-grid-column
                field="class_name"
                title="班级名称"
                align="center"
              ></tiny-grid-column>
              <tiny-grid-column
                field="course_name"
                title="课程名称"
                align="center"
              ></tiny-grid-column>
              <tiny-grid-column
                field="user_name"
                title="教师姓名"
                align="center"
              ></tiny-grid-column>
              <tiny-grid-column
                field="question_set_name"
                title="问题集名称"
                align="center"
              ></tiny-grid-column>
              <tiny-grid-column
                field="title"
                title="问题标题"
                align="left"
              ></tiny-grid-column>
              <tiny-grid-column
                field="options"
                title="选项"
                align="left"
              ></tiny-grid-column>
              <tiny-grid-column
                field="qty"
                title="选择数量"
                align="center"
              ></tiny-grid-column>
            </tiny-grid>
          </div>
        </div>
      </tiny-fullscreen>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, toRefs, onMounted } from 'vue';
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
  } from '@opentiny/vue';
  import { request } from '@/utils/request';
  import Breadcrumb from '@/components/breadcrumb/index.vue';
  import { queryEvaluationTaskList } from '@/api/evaluationTask';
  import { queryClassList } from '@/api/class';
  import { queryCourseList } from '@/api/course';
  import { getAdminUserList } from '@/api/user';

  const taskLoading = ref(false);
  const taskDataList = ref([]);
  const classesLoading = ref(false);
  const classesDataList = ref([]);
  const courseLoading = ref(false);
  const courseDataList = ref([]);
  const adminUserLoading = ref(false);
  const adminUserDataList = ref([]);

  // 初始化请求数据
  interface FilterOptions {
    taskId: number;
    classId: number;
    courseId: number;
    userId: number;
  }

  // 加载效果
  const state = reactive<{
    loading: boolean;
    filterOptions: FilterOptions;
  }>({
    loading: false,
    filterOptions: {
      taskId: undefined,
      classId: undefined,
      courseId: undefined,
      userId: undefined,
    },
  });

  const pagerConfig = reactive({
    component: TinyPager,
    attrs: {
      currentPage: 1,
      pageSize: 10,
      pageSizes: [10, 20, 50],
      total: 10,
      layout: 'total, prev, pager, next, jumper, sizes',
    },
  });

  let tableData = ref([]);
  const taskGrid = ref();
  const { loading, filterOptions } = toRefs(state);

  // 请求数据接口方法
  async function fetchData(
    params: any = {
      page: 1,
      limit: 10,
    },
  ) {
    const { ...rest } = filterOptions.value;
    const queryParmas = {
      ...params,
      ...rest,
    };

    // 转换数字类型，确保ID参数为数字
    Object.keys(queryParmas).forEach((key) => {
      if (['taskId', 'classId', 'courseId', 'userId'].includes(key)) {
        queryParmas[key] = Number(queryParmas[key]) || 0;
      }
    });

    state.loading = true;
    try {
      // 获取评教详情统计列表
      const res: any = await request.get(
        {
          url: '/manager/api/statistics/list',
          params: queryParmas,
        },
        { isTransformResponse: false },
      );

      if (res.code !== 200) {
        Modal.message({
          message: res.msg || '获取评教详情统计失败',
          status: 'error',
        });
        return {
          result: [],
          page: { total: 0 },
        };
      }

      tableData.value = res.data.list;
      return {
        result: res.data.list,
        page: { total: res.data.total },
      };
    } catch (error: any) {
      Modal.message({
        message: error.message || '获取评教详情统计失败',
        status: 'error',
      });
      return {
        result: [],
        page: { total: 0 },
      };
    } finally {
      state.loading = false;
    }
  }

  
  async function loadTaskData(query: string) {
    console.log('query:', query);
    taskLoading.value = true;
    const queryParams = {
      page: 1,
      limit: 10,
      keyword: query?.trim() || '',
    };
    try {
      const res: any = await queryEvaluationTaskList(queryParams);
      console.log('res:', res);
      taskDataList.value = res.data;
    } finally {
      taskLoading.value = false;
    }
  }
  async function loadClassesData(query: string) {
    console.log('query:', query);
    classesLoading.value = true;
    const queryParams = {
      page: 1,
      limit: 10,
      keyword: query?.trim() || '',
    };
    try {
      const res: any = await queryClassList(queryParams);
      console.log('res:', res);
      if (res.code === 200) {
        classesDataList.value = res.data;
      }
    } finally {
      classesLoading.value = false;
    }
  }
  async function loadCourseData(query: string) {
    console.log('query:', query);
    courseLoading.value = true;
    const queryParams = {
      page: 1,
      limit: 10,
      keyword: query?.trim() || '',
    };
    try {
      const res: any = await queryCourseList(queryParams);
      console.log('res:', res);
      if (res.code === 200) {
        courseDataList.value = res.data;
      }
    } finally {
      courseLoading.value = false;
    }
  }
  async function loadAdminUserData(query: string) {
    console.log('query:', query);
    adminUserLoading.value = true;
    const queryParams = {
      page: 1,
      limit: 10,
      isPage: 0,
      keyword: query?.trim() || '',
      account: ''
    };
    try {
      const res: any = await getAdminUserList(queryParams);
      console.log('res:', res);
      if (res.code === 200) {
        adminUserDataList.value = res.data.list;
      }
    } finally {
      adminUserLoading.value = false;
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

  // 导出评教详情统计为CSV
  const exportStats = async () => {
    const { ...rest } = filterOptions.value;
    const queryParmas = {
      ...rest,
    };

    // 转换数字类型，确保ID参数为数字
    Object.keys(queryParmas).forEach((key) => {
      if (['taskId', 'classId', 'courseId', 'userId'].includes(key)) {
        queryParmas[key] = Number(queryParmas[key]) || 0;
      }
    });

    try {
      // 导出评教详情统计为CSV
      const res = await request.get(
        {
          url: '/manager/api/statistics/export',
          params: queryParmas,
          responseType: 'blob',
        },
        { isTransformResponse: false },
      );

      // 处理CSV文件下载
      const blob = new Blob([res], { type: 'text/csv;charset=utf-8;' });
      const link = document.createElement('a');
      const url = URL.createObjectURL(blob);
      link.setAttribute('href', url);
      link.setAttribute('download', `评教详情统计_${new Date().getTime()}.csv`);
      link.style.visibility = 'hidden';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);

      Modal.message({
        message: '导出成功',
        status: 'success',
      });
    } catch (error: any) {
      Modal.message({
        message: error.message || '导出失败',
        status: 'error',
      });
    }
  };

  // form的button
  function reloadGrid() {
    taskGrid?.value.handleFetch('reload');
  }

  function handleFormReset() {
    state.filterOptions = {
      taskId: 0,
      classId: 0,
      courseId: 0,
      userId: 0,
    };
    reloadGrid();
  }

  // 全屏缩放设置
  const fullscreen = ref(false);
  const toggle = () => {
    fullscreen.value = !fullscreen.value;
  };

  onMounted(() => {
    // reloadGrid();
    loadTaskData('');
    loadClassesData('');
    loadCourseData('');
    loadAdminUserData('');
  });
</script>

<style scoped lang="less">
  @import '@/assets/style/page-table.less';
</style>

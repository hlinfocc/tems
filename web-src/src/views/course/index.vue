<template>
  <div class="container-list">
    <Breadcrumb :items="['课程管理', '课程列表']" :is-i18n="false" />
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
            <tiny-form-item label="课程名称">
              <tiny-input
                v-model="filterOptions.courseName"
                placeholder="请输入课程名称"
              ></tiny-input>
            </tiny-form-item>
          </tiny-col>
          <tiny-col :span="4" label-width="100px">
            <tiny-form-item label="状态">
              <tiny-select
                v-model="filterOptions.status"
                placeholder="请选择状态"
              >
                <tiny-option label="全部" value="" />
                <tiny-option label="正常" value="0" />
                <tiny-option label="禁用" value="1" />
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
                        @click="() => (editorVisible = true)"
                      >
                        添加
                      </tiny-button>
                      <tiny-button
                        type="info"
                        class="batch-btn"
                        @click="() => (barchAddVisible = true)"
                      >
                        批量添加
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
                field="courseName"
                title="课程名称"
                align="center"
              ></tiny-grid-column>
              <tiny-grid-column field="status" title="状态" align="center">
                <template #default="{ row }">
                  <tiny-tag :type="row.status === 0 ? 'success' : 'danger'">
                    {{ row.status === 0 ? '正常' : '禁用' }}
                  </tiny-tag>
                </template>
              </tiny-grid-column>
              <tiny-grid-column
                field="updatedAt"
                title="更新时间"
                align="center"
              ></tiny-grid-column>
              <tiny-grid-column title="操作" align="center">
                <template v-slot="data">
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
    <CourseEditor
      v-model:visible="editorVisible"
      :course-data="selectedCourse"
      @success="addCallback"
    />
    <batch-add  v-model:visible="barchAddVisible" @success="addCallback" />
  </div>
</template>

<script setup lang="ts">
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
  import { queryCourseList, deleteCourse } from '@/api/course';
  import type { CourseInfo, QueryParmas } from '@/api/course';
  import Breadcrumb from '@/components/breadcrumb/index.vue';
  import CourseEditor from './components/editor.vue';
  import batchAdd from './components/batchAdd.vue';

  const editorVisible = ref(false);
  const barchAddVisible = ref(false);
  const selectedCourse = ref<any>({});

  // 初始化请求数据
  interface FilterOptions {
    courseName: string;
    status: string;
  }

  // 加载效果
  const state = reactive<{
    loading: boolean;
    filterOptions: FilterOptions;
  }>({
    loading: false,
    filterOptions: {
      courseName: '',
      status: '',
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
    params: QueryParmas = {
      page: 1,
      limit: 10,
      keyword: '',
      status: '',
    },
  ) {
    console.log('filterOptions.value:', filterOptions.value);

    const { ...rest } = filterOptions.value;
    const queryParmas = {
      ...params,
      ...rest,
    };

    state.loading = true;
    try {
      const res: any = await queryCourseList(queryParmas);

      if (res.code !== 200) {
        Modal.message({
          message: res.msg || '获取课程列表失败',
          status: 'error',
        });
        return {
          result: [],
          page: { total: 0 },
        };
      }
      tableData.value = res.data;
      return {
        result: res.data,
        page: { total: res.count },
      };
    } catch (error: any) {
      Modal.message({
        message: error.message || '获取课程列表失败',
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

  const addCallback = () => {
    reloadGrid();
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
    Modal.confirm('您确定要删除该课程吗？').then(async (rs: any) => {
      try {
        const res: any = await deleteCourse(id);
        if (res.code === 200) {
          reloadGrid();
          Modal.message({
            message: res.msg,
            status: 'success',
          });
        } else {
          Modal.message({
            message: res.msg || '删除课程失败',
            status: 'error',
          });
        }
      } catch (error: any) {
        Modal.message({
          message: error.message || '删除课程失败',
          status: 'error',
        });
      }
    });
  };

  const handleEdit = (e: any) => {
    editorVisible.value = true;
    console.log('course data:', toRaw(e));
    selectedCourse.value = e;
  };

  // form的button
  function reloadGrid() {
    taskGrid?.value.handleFetch('reload');
    // fetchData();
  }

  function handleFormReset() {
    state.filterOptions = {
      courseName: '',
      status: '',
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
  });
</script>

<style scoped lang="less">
  @import '@/assets/style/page-table.less';
  .btn {
    display: flex;
    justify-content: start;
    width: 100%;
    position: relative;
  }
  .batch-btn {
    width: 110px !important;
  }
  .screen {
    position: absolute;
    right: 7px;
  }
</style>

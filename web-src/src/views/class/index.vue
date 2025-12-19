<template>
  <div class="container-list">
    <Breadcrumb :items="['班级管理', '班级列表']" :is-i18n="false" />
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
            <tiny-form-item label="班级名称">
              <tiny-input
                v-model="filterOptions.keyword"
                placeholder="请输入班级名称"
              ></tiny-input>
            </tiny-form-item>
          </tiny-col>
          <tiny-col :span="4" label-width="100px">
            <tiny-form-item label="专业">
              <tiny-input
                v-model="filterOptions.major"
                placeholder="请输入专业"
              ></tiny-input>
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
              ref="gridRef"
              :fetch-data="fetchDataOption"
              :pager="pagerConfig"
              :loading="loading"
              seq-serial
              size="medium"
              :auto-resize="true"
              :auto-load="false"
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
                field="className"
                title="班级名称"
                align="center"
              >
              </tiny-grid-column>
              <tiny-grid-column field="majorName" title="专业" align="center">
              </tiny-grid-column>
              <tiny-grid-column field="grade" title="年级" align="center">
              </tiny-grid-column>
              <tiny-grid-column
                field="headTeacherName"
                title="班主任"
                align="center"
              >
              </tiny-grid-column>

              <tiny-grid-column
                field="updatedAt"
                title="更新时间"
                align="center"
              >
              </tiny-grid-column>
              <tiny-grid-column title="操作" align="center">
                <template v-slot="data">
                  <a
                    class="operation-item"
                    @click="handleCourseManage(data.row)"
                  >
                    课程管理
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
      :class-data="selectedClass"
      @success="addCallback"
    />
    <CourseManager
      v-model:visible="courseManagerVisible"
      :class-data="selectedClassForCourse"
    />
  </div>
</template>

<script lang="ts" setup>
  import type { QueryParmas } from '@/api/class';
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
    Pager as TinyPager,
    Fullscreen as TinyFullscreen,
    Modal,
    TinyTag,
  } from '@opentiny/vue';
  import { queryClassList, deleteClass } from '@/api/class';
  import Breadcrumb from '@/components/breadcrumb/index.vue';
  import Editor from './components/editor.vue';
  import CourseManager from './components/courseManager.vue';

  const editorVisible = ref(false);
  const selectedClass = ref<any>({});
  const courseManagerVisible = ref(false);
  const selectedClassForCourse = ref<any>({});

  // 初始化请求数据
  interface FilterOptions {
    keyword: string;
    status: string;
    major: string;
    grade: string;
  }

  // 加载效果
  const state = reactive<{
    loading: boolean;
    filterOptions: FilterOptions;
  }>({
    loading: false,
    filterOptions: {} as FilterOptions,
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
  const gridRef = ref();
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
    console.log('params:', params);
    const { ...rest } = filterOptions.value;
    const queryParmas = {
      ...params,
      ...rest,
    };

    state.loading = true;
    try {
      const res: any = await queryClassList(queryParmas);
      console.log('res:', res);
      if (res.code !== 200) {
        Modal.message({
          message: res.msg || '获取班级列表失败',
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
        message: error.message || '获取班级列表失败',
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
      console.log('currentPage:', currentPage);
      return fetchData({
        page: currentPage,
        limit: pageSize,
      });
    },
  });

  const handleDelete = (id: number) => {
    Modal.confirm('您确定要删除该班级吗？').then(async (rs: any) => {
      if (rs !== 'confirm') {
        return;
      }
      try {
        const res: any = await deleteClass(id);
        if (res.code === 200) {
          reloadGrid();
          Modal.message({
            message: res.msg,
            status: 'success',
          });
        } else {
          Modal.message({
            message: res.msg || '删除班级失败',
            status: 'error',
          });
        }
      } catch (error: any) {
        Modal.message({
          message: error.message || '删除班级失败',
          status: 'error',
        });
      }
    });
  };

  const handleEdit = (e: any) => {
    editorVisible.value = true;
    console.log('class data:', toRaw(e));
    selectedClass.value = e;
  };

  const handleCourseManage = (e: any) => {
    courseManagerVisible.value = true;
    selectedClassForCourse.value = e;
  };

  // form的button
  function reloadGrid() {
    gridRef?.value.handleFetch('reload');
    // fetchData();
  }

  function handleFormReset() {
    state.filterOptions = {} as FilterOptions;
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

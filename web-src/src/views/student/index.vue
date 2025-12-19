<template>
  <div class="container-list">
    <Breadcrumb :items="['学生管理', '学生列表']" :isI18n="false" />
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
            <tiny-form-item label="姓名">
              <tiny-input
                v-model="filterOptions.studentName"
                placeholder="请输入姓名"
                clearable
              ></tiny-input>
            </tiny-form-item>
          </tiny-col>
          <tiny-col :span="4" label-width="100px">
            <tiny-form-item label="学号">
              <tiny-input
                v-model="filterOptions.studentID"
                placeholder="请输入学号"
                clearable
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
                            editorVisible = true;
                          }
                        "
                      >
                        添加
                      </tiny-button>
                      <tiny-button
                        type="info"
                        class="batch-btn"
                        @click="
                          () => {
                            batchAddVisible = true;
                          }
                        "
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
              <tiny-grid-column field="studentName" title="姓名" align="center">
              </tiny-grid-column>
              <tiny-grid-column field="studentID" title="学号" align="center">
              </tiny-grid-column>
              <tiny-grid-column field="phone" title="手机号" align="center">
              </tiny-grid-column>
              <tiny-grid-column field="className" title="班级" align="center">
              </tiny-grid-column>
              <tiny-grid-column field="isBound" title="绑定状态" align="center">
                <template #default="{ row }">
                  <tiny-tag :type="row.isBound ? 'success' : 'danger'">
                    {{ row.isBound ? '已绑定' : '未绑定' }}
                  </tiny-tag>
                </template>
              </tiny-grid-column>
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
    <editor
      v-model:visible="editorVisible"
      :student-data="selectedStudent"
      @success="addCallback"
    />
    <batch-add v-model:visible="batchAddVisible" @success="addCallback" />
  </div>
</template>

<script lang="ts" setup>
  import type { QueryParmas } from '@/api/student';
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
  import { queryStudentList, deleteStudent } from '@/api/student';
  import Breadcrumb from '@/components/breadcrumb/index.vue';
  import Editor from './components/editor.vue';
  import BatchAdd from './components/batchAdd.vue';

  const editorVisible = ref(false);
  const batchAddVisible = ref(false);
  const selectedStudent = ref<any>({});

  // 初始化请求数据
  interface FilterOptions {
    studentName: string;
    studentID: string;
    page: number;
    pageSize: number;
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
  const taskGrid = ref();
  const { loading, filterOptions } = toRefs(state);

  const statusOptions = [
    {
      value: '0',
      label: '正常',
    },
    {
      value: '1',
      label: '禁用',
    },
  ];

  // 请求数据接口方法
  async function fetchData(
    params: QueryParmas = {
      page: 1,
      pageSize: 10,
      studentName: '',
      studentID: '',
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
      const res: any = await queryStudentList(queryParmas);

      tableData.value = res.list;
      return {
        result: res.list,
        page: { total: res.total },
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
    Modal.confirm('您确定要删除该学生吗？').then((rs: any) => {
      deleteStudent(id).then((res: any) => {
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
    console.log('student data:', toRaw(e));
    selectedStudent.value = e;
  };

  // form的button
  function reloadGrid() {
    taskGrid?.value.handleFetch('reload');
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

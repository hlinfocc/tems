<template>
  <div class="container-list">
    <Breadcrumb :items="['评教管理', '评教问题集']" :is-i18n="false" />
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
            <tiny-form-item label="问题集名称">
              <tiny-input
                v-model="filterOptions.keyword"
                placeholder="请输入问题集名称"
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
              <tiny-grid-column field="name" title="问题集名称" align="center">
              </tiny-grid-column>
              <tiny-grid-column field="remark" title="备注" align="center">
              </tiny-grid-column>
              <tiny-grid-column field="year" title="年份" align="center">
              </tiny-grid-column>
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
              >
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
      :question-set-data="selectedQuestionSet"
      @success="addCallback"
    />
  </div>
</template>

<script lang="ts" setup>
  import type { QueryParmas } from '@/api/evaluationQuestionSet';
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
  import {
    queryEvaluationQuestionSetList,
    deleteEvaluationQuestionSet,
  } from '@/api/evaluationQuestionSet';
  import Breadcrumb from '@/components/breadcrumb/index.vue';
  import Editor from './components/editor.vue';

  const editorVisible = ref(false);
  const selectedQuestionSet = ref<any>({});

  // 初始化请求数据
  interface FilterOptions {
    keyword: string;
    status: string;
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
      const res: any = await queryEvaluationQuestionSetList(queryParmas);
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
    Modal.confirm('您确定要删除该评教问题集吗？').then((rs: any) => {
      deleteEvaluationQuestionSet(id).then((res: any) => {
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
    console.log('question set data:', toRaw(e));
    selectedQuestionSet.value = e;
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
</style>

<template>
  <div class="container-list">
    <Breadcrumb :items="['评教管理', '评教问题详情']" :is-i18n="false" />
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
            <tiny-form-item label="问题标题">
              <tiny-input
                v-model="filterOptions.keyword"
                placeholder="请输入问题标题"
              ></tiny-input>
            </tiny-form-item>
          </tiny-col>
          <tiny-col :span="4" label-width="100px">
            <tiny-form-item label="问题集">
              <tiny-select
                v-model="filterOptions.questionSetId"
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
              <tiny-grid-column field="title" title="问题标题" align="left">
              </tiny-grid-column>
              <tiny-grid-column
                field="questionType"
                title="问题类型"
                align="center"
                width="120"
              >
                <template #default="{ row }">
                  <tiny-tag v-if="row.questionType == 1">单选</tiny-tag>
                  <tiny-tag v-if="row.questionType == 2">多选</tiny-tag>
                </template>
              </tiny-grid-column>
              <tiny-grid-column field="remark" title="备注" align="left">
              </tiny-grid-column>
              <tiny-grid-column
                field="updatedAt"
                title="更新时间"
                align="center"
                width="180"
              >
              </tiny-grid-column>
              <tiny-grid-column title="操作" align="center" width="160">
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
      :question-data="selectedQuestion"
      @success="addCallback"
    />
  </div>
</template>

<script lang="ts" setup>
  import type { QueryParams } from '@/api/evaluation-question-detail';
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
    Option as TinyOption,
    Modal,
    TinyTag,
  } from '@opentiny/vue';
  import {
    queryQuestionDetailList,
    deleteQuestionDetail,
  } from '@/api/evaluation-question-detail';
  import { queryEvaluationQuestionSetList } from '@/api/evaluationQuestionSet';
  import Breadcrumb from '@/components/breadcrumb/index.vue';
  import Editor from './components/editor.vue';

  const editorVisible = ref(false);
  const selectedQuestion = ref<any>({});
  const questionSetListData = ref([]);
  // 初始化请求数据
  interface FilterOptions {
    keyword: string;
    questionSetId: string;
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
    params: QueryParams = {
      page: 1,
      limit: 10,
      keyword: '',
      questionSetId: '',
    },
  ) {
    console.log('filterOptions.value:', filterOptions.value);

    const { ...rest } = filterOptions.value;
    const queryParams = {
      ...params,
      ...rest,
    };

    state.loading = true;
    try {
      const res: any = await queryQuestionDetailList(queryParams);
      if (res.code === 200) {
        tableData.value = res.data;
      }
      return {
        result: res.data,
        page: { total: res.count },
      };
    } finally {
      state.loading = false;
    }
  }

  async function fetchEQSetsData() {
    let params: any = {
      page: 1,
      limit: 10,
      isPage: 1,
      keyword: '',
      status: '',
    };
    const res: any = await queryEvaluationQuestionSetList(params);
    console.log('res:', res);
    if (res.code === 200) {
      questionSetListData.value = res.data;
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
    Modal.confirm('您确定要删除该评教问题详情吗？').then((rs: any) => {
      deleteQuestionDetail(id).then((res: any) => {
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
    console.log('question data:', toRaw(e));
    selectedQuestion.value = e;
  };

  // form的button
  function reloadGrid() {
    taskGrid?.value.handleFetch('reload');
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
    fetchEQSetsData();
    reloadGrid();
  });
</script>

<style scoped lang="less">
  @import '@/assets/style/page-table.less';
</style>

<template>
  <div class="container-list">
    <Breadcrumb :items="['轮播图管理', '轮播图列表']" :is-i18n="false" />
    <div class="contain">
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
              <tiny-grid-column field="title" title="轮播图标题" align="center">
              </tiny-grid-column>
              <tiny-grid-column
                field="image_url"
                title="图片URL"
                align="center"
              >
                <template v-slot="data">
                  <div class="image-preview">
                    <img
                      :src="data.row.image_url"
                      alt="轮播图"
                      class="preview-img"
                      @click="previewImage(data.row.image_url)"
                    />
                  </div>
                </template>
              </tiny-grid-column>
              <tiny-grid-column field="sort" title="排序" align="center">
              </tiny-grid-column>
              <tiny-grid-column
                field="is_visible"
                title="是否可见"
                align="center"
              >
                <template v-slot="data">
                  <tiny-switch
                    v-model="data.row.is_visible"
                    @change="handleVisibilityChange(data.row)"
                  ></tiny-switch>
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
      :banner-data="selectedBanner"
      @success="addCallback"
    />
  </div>
</template>

<script lang="ts" setup>
  import type { QueryParams } from '@/api/banner';
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
    Fullscreen as TinyFullscreen,
    Switch as TinySwitch,
    Pager as TinyPager,
    Modal,
  } from '@opentiny/vue';
  import {
    queryBannerList,
    deleteBanner,
    updateBannerVisibility,
  } from '@/api/banner';
  import Breadcrumb from '@/components/breadcrumb/index.vue';
  import Editor from './components/editor.vue';

  const editorVisible = ref(false);
  const selectedBanner = ref<any>({});

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
  const gridRef = ref();
  const { loading, filterOptions } = toRefs(state);

  // 请求数据接口方法
  async function fetchData(
    params: QueryParams = {
      page: 1,
      limit: 10,
      keyword: '',
    },
  ) {
    console.log('filterOptions.value:', filterOptions.value);
    console.log('params:', params);
    const { ...rest } = filterOptions.value;
    const queryParams = {
      ...params,
      ...rest,
    };

    state.loading = true;
    try {
      const res: any = await queryBannerList(queryParams);
      console.log('res:', res);
      // if (res.code !== 200) {
      //   Modal.message({
      //     message: res.msg || '获取轮播图列表失败',
      //     status: 'error',
      //   });
      //   return {
      //     result: [],
      //     page: { total: 0 },
      //   };
      // }
      tableData.value = res.list;
      return {
        result: res.list,
        page: { total: res.total },
      };
    } catch (error: any) {
      Modal.message({
        message: error.message || '获取轮播图列表失败',
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
    Modal.confirm('您确定要删除该轮播图吗？').then(async (rs: any) => {
      if (rs !== 'confirm') {
        return;
      }
      try {
        const res: any = await deleteBanner(id);
        if (res.code === 200) {
          reloadGrid();
          Modal.message({
            message: res.msg,
            status: 'success',
          });
        } else {
          Modal.message({
            message: res.msg || '删除轮播图失败',
            status: 'error',
          });
        }
      } catch (error: any) {
        Modal.message({
          message: error.message || '删除轮播图失败',
          status: 'error',
        });
      }
    });
  };

  const handleEdit = (e: any) => {
    editorVisible.value = true;
    console.log('banner data:', toRaw(e));
    selectedBanner.value = e;
  };

  // 处理可见性变化
  const handleVisibilityChange = async (row: any) => {
    try {
      const res: any = await updateBannerVisibility({
        id: row.id,
        is_visible: row.is_visible,
      });
      if (res.code === 200) {
        Modal.message({
          message: res.msg,
          status: 'success',
        });
      } else {
        // 恢复原来的状态
        row.is_visible = !row.is_visible;
        Modal.message({
          message: res.msg || '更新可见性失败',
          status: 'error',
        });
      }
    } catch (error: any) {
      // 恢复原来的状态
      row.is_visible = !row.is_visible;
      Modal.message({
        message: error.message || '更新可见性失败',
        status: 'error',
      });
    }
  };

  // 预览图片
  const previewImage = (url: string) => {
    // 这里可以实现图片预览功能
    window.open(url, '_blank');
  };

  // form的button
  function reloadGrid() {
    gridRef?.value.handleFetch('reload');
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

  .image-preview {
    display: flex;
    justify-content: center;
    align-items: center;

    .preview-img {
      width: 100px;
      height: 50px;
      object-fit: cover;
      cursor: pointer;
      border-radius: 4px;
      border: 1px solid #e8e8e8;
    }
  }
</style>

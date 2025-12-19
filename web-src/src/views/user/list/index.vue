<template>
  <div class="container-list">
    <Breadcrumb :items="['用户管理', '管理员列表']" :is-i18n="false" />
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
                v-model="filterOptions.name"
                placeholder="请输入姓名"
              ></tiny-input>
            </tiny-form-item>
          </tiny-col>
          <tiny-col :span="4" label-width="100px">
            <tiny-form-item label="用户名">
              <tiny-input
                v-model="filterOptions.username"
                placeholder="请输入用户名"
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
                <tiny-option label="启用" value="0" />
                <tiny-option label="禁用" value="1" />
              </tiny-select>
            </tiny-form-item>
          </tiny-col>
          <tiny-col :span="4" label-width="100px">
            <tiny-form-item label="用户类型">
              <tiny-select
                v-model="filterOptions.user_type"
                placeholder="请选择用户类型"
              >
                <tiny-option label="全部" value="" />
                <tiny-option label="管理员" value="0" />
                <tiny-option label="老师" value="1" />
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
                        @click="() => {editorVisible = true, isUserEdit = false}"
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
                field="name"
                title="姓名"
                align="center"
              ></tiny-grid-column>
              <tiny-grid-column
                field="username"
                title="用户名"
                align="center"
              ></tiny-grid-column>
              <tiny-grid-column field="status" title="状态" align="center">
                <template #default="{ row }">
                  <tiny-tag :type="row.status === 0 ? 'success' : 'danger'">
                    {{ row.status === 0 ? '启用' : '禁用' }}
                  </tiny-tag>
                </template>
              </tiny-grid-column>
              <tiny-grid-column
                field="user_type"
                title="用户类型"
                align="center"
              >
                <template #default="{ row }">
                  <tiny-tag :type="row.user_type === 0 ? 'info' : 'warning'">
                    {{ row.user_type === 0 ? '管理员' : '老师' }}
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
                  <a
                    class="operation-item"
                    @click="handleResetPassword(data.row.id)"
                  >
                    重置密码
                  </a>
                </template>
              </tiny-grid-column>
            </tiny-grid>
          </div>
        </div>
      </tiny-fullscreen>
    </div>
    <AdminUserEditor
      v-model:visible="editorVisible"
      :user-data="selectedUser"
      :is-edit="isUserEdit"
      @success="addCallback"
    />
    <tiny-dialog-box
      v-model:visible="resetPasswordVisible"
      title="重置密码"
      width="400px"
      :mask-closable="false"
    >
      <tiny-form
        ref="resetPasswordForm"
        :model="resetPasswordData"
        :rules="resetPasswordRules"
        label-position="top"
      >
        <tiny-form-item label="新密码" prop="password">
          <tiny-input
            v-model="resetPasswordData.password"
            type="password"
            placeholder="请输入新密码"
            :maxlength="20"
            show-word-limit
          ></tiny-input>
        </tiny-form-item>
      </tiny-form>
      <template #footer>
        <tiny-button @click="resetPasswordVisible = false">取消</tiny-button>
        <tiny-button
          type="info"
          @click="handleSubmitResetPassword"
          :loading="resetPasswordLoading"
          >确定</tiny-button
        >
      </template>
    </tiny-dialog-box>
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
    TinyDialogBox,
  } from '@opentiny/vue';
  import {
    getAdminUserList,
    deleteAdminUser,
    resetAdminUserPassword,
  } from '@/api/user';
  import type {
    AdminUserListItem,
    AdminUserListParams,
    ResetPasswordRequest,
  } from '@/api/user';
  import Breadcrumb from '@/components/breadcrumb/index.vue';
  import AdminUserEditor from './components/editor.vue';

  const editorVisible = ref(false);
  const resetPasswordVisible = ref(false);
  const selectedUser = ref<any>({});
  const resetPasswordUserId = ref<number>(0);
  const resetPasswordForm = ref();
  const resetPasswordLoading = ref(false);
  const isUserEdit = ref(false);

  // 初始化请求数据
  interface FilterOptions {
    name: string;
    username: string;
    status: string;
    user_type: string;
  }

  // 加载效果
  const state = reactive<{
    loading: boolean;
    filterOptions: FilterOptions;
  }>({
    loading: false,
    filterOptions: {
      name: '',
      username: '',
      status: '',
      user_type: '',
    },
  });

  // 重置密码表单数据
  const resetPasswordData = reactive<ResetPasswordRequest>({
    password: '',
  });

  // 重置密码表单规则
  const resetPasswordRules = {
    password: [
      { required: true, message: '请输入新密码', trigger: 'blur' },
      {
        min: 6,
        max: 20,
        message: '密码长度在 6 到 20 个字符',
        trigger: 'blur',
      },
    ],
  };

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
    params: AdminUserListParams = {
      page: 1,
      limit: 10,
      isPage: 0,
      keyword: '',
      account: '',
      status: undefined,
      user_type: undefined,
    },
  ) {
    const queryParams: AdminUserListParams = {
      ...params,
      keyword: filterOptions.value.name,
      account: filterOptions.value.username,
      isPage: 0,
    };

    if (filterOptions.value.status) {
      queryParams.status = parseInt(filterOptions.value.status as string, 10);
    }

    if (filterOptions.value.user_type) {
      queryParams.user_type = parseInt(filterOptions.value.user_type as string, 10);
    }

    state.loading = true;
    try {
      const res: any = await getAdminUserList(queryParams);

      if (res.code !== 200) {
        Modal.message({
          message: res.msg || '获取管理员列表失败',
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
        message: error.message || '获取管理员列表失败',
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
    Modal.confirm('您确定要删除该管理员吗？').then(async (rs: any) => {
      try {
        const res: any = await deleteAdminUser(id);
        if (res.code === 200) {
          reloadGrid();
          Modal.message({
            message: res.msg,
            status: 'success',
          });
        } else {
          Modal.message({
            message: res.msg || '删除管理员失败',
            status: 'error',
          });
        }
      } catch (error: any) {
        Modal.message({
          message: error.message || '删除管理员失败',
          status: 'error',
        });
      }
    });
  };

  const handleEdit = (e: any) => {
    editorVisible.value = true;
    selectedUser.value = e;
    isUserEdit.value = true;
  };

  const handleResetPassword = (id: number) => {
    resetPasswordUserId.value = id;
    resetPasswordVisible.value = true;
    resetPasswordData.password = '';
    if (resetPasswordForm.value) {
      resetPasswordForm.value.resetFields();
    }
  };

  const handleSubmitResetPassword = async () => {
    if (!resetPasswordForm.value) return;
    resetPasswordForm.value.validate(async (valid: boolean) => {
      if (valid) {
        try {
          resetPasswordLoading.value = true;
          const res: any = await resetAdminUserPassword(
            resetPasswordUserId.value,
            resetPasswordData,
          );
          if (res.code === 200) {
            resetPasswordVisible.value = false;
            Modal.message({
              message: res.msg,
              status: 'success',
            });
          } else {
            Modal.message({
              message: res.msg || '重置密码失败',
              status: 'error',
            });
          }
        } catch (error: any) {
          Modal.message({
            message: error.message || '重置密码失败',
            status: 'error',
          });
        } finally {
          resetPasswordLoading.value = false;
        }
      }
    });
  };

  // form的button
  function reloadGrid() {
    taskGrid?.value.handleFetch('reload');
  }

  function handleFormReset() {
    state.filterOptions = {
      name: '',
      username: '',
      status: '',
      user_type: '',
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
</style>

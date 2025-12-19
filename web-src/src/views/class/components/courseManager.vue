<template>
  <tiny-drawer
    v-model:visible="visible"
    title="班级课程管理"
    :width="`85%`"
    :mask-closable="false"
    @close="handleClose"
  >
    <div class="detail-container">
      <div class="task-info">
        <h3>{{ classData.className }} - 课程列表</h3>
        <div class="task-meta">
          <span>专业：{{ classData.majorName }}</span>
          <span>年级：{{ classData.grade }}</span>
          <span>班主任：{{ classData.headTeacherName }}</span>
        </div>
      </div>

      <div class="detail-actions">
        <tiny-button type="info" @click="handleAddCourse">添加</tiny-button>
        <!-- <tiny-button
          type="danger"
          @click="handleBatchDelete"
          :disabled="selectedRows.length === 0"
          >批量删除</tiny-button
        > -->
      </div>

      <tiny-grid
        ref="courseGrid"
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
        <tiny-grid-column field="courseName" title="课程名称" align="center">
        </tiny-grid-column>
        <tiny-grid-column field="teacherName" title="任课教师" align="center">
        </tiny-grid-column>
        <tiny-grid-column field="updatedAt" title="更新时间" align="center">
        </tiny-grid-column>
        <tiny-grid-column title="操作" align="center">
          <template v-slot="data">
            <a class="operation-item" @click="handleDeleteCourse(data.row.id)">
              删除
            </a>
          </template>
        </tiny-grid-column>
      </tiny-grid>
    </div>

    <!-- 课程编辑子抽屉 -->
    <CourseEditor
      v-model:visible="editorVisible"
      :course-data="selectedCourse"
      :class-id="classData.id"
      :class-name="classData.className"
      @success="handleCourseSuccess"
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
  import type { ClassInfo } from '@/api/class';
  import type {
    ClassCourseInfo,
    QueryParams,
  } from '@/api/classCourse';
  import {
    getClassCourseList,
    deleteClassCourse,
  } from '@/api/classCourse';
  import CourseEditor from './courseEditor.vue';

  const visible = ref<boolean>(false);
  const editorVisible = ref<boolean>(false);
  const selectedCourse = ref<any>({});
  const selectedRows = ref<any[]>([]);

  const props = defineProps<{
    visible: boolean;
    classData: ClassInfo;
  }>();

  const emit = defineEmits(['update:visible', 'success']);
  const propsVisible = toRefs(props).visible;
  const propsClassData = toRefs(props).classData;

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
  const courseGrid = ref();
  const classData = ref<ClassInfo>({} as ClassInfo);

  // 请求数据接口方法
  async function fetchData(
    params: QueryParams = {
      page: 1,
      limit: 10,
      classId: undefined,
    },
  ) {
    if (!classData.value.id) return { result: [], page: { total: 0 } };

    const queryParams = {
      ...params,
      classId: classData.value.id,
    };

    state.loading = true;
    try {
      const res: any = await getClassCourseList(queryParams);
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

  // 添加课程
  const handleAddCourse = () => {
    selectedCourse.value = { classId: classData.value.id };
    editorVisible.value = true;
  };

  // 编辑课程
  const handleEditCourse = (row: ClassCourseInfo) => {
    selectedCourse.value = { ...row };
    editorVisible.value = true;
  };

  // 删除课程
  const handleDeleteCourse = (id: number) => {
    Modal.confirm('您确定要删除该班级课程吗？').then((rs: any) => {
      deleteClassCourse(id).then((res: any) => {
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
        message: '请选择要删除的课程',
        status: 'warning',
      });
      return;
    }

    Modal.confirm(
      `您确定要删除选中的${selectedRows.value.length}条班级课程吗？`,
    ).then(async (rs: any) => {
      try {
        // 使用Promise.all并行执行删除操作
        await Promise.all(
          selectedRows.value.map((row) => deleteClassCourse(row.id)),
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

  // 课程操作成功回调
  const handleCourseSuccess = () => {
    reloadGrid();
    emit('success');
  };

  // 刷新表格
  const reloadGrid = () => {
    courseGrid?.value.handleFetch('reload');
    fetchData();
  };

  // 关闭抽屉
  const handleClose = () => {
    visible.value = false;
    emit('update:visible', false);
    selectedRows.value = [];
  };

  // 监听班级数据变化
  watch(
    () => propsClassData.value,
    (val) => {
      if (val && val.id) {
        classData.value = { ...val };
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
    if (newValue && classData.value.id) {
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
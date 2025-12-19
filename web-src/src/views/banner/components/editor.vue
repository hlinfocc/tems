<template>
  <tiny-drawer
    v-model:visible="visible"
    :title="isEdit ? '编辑轮播图' : '新增轮播图'"
    :width="`40%`"
    show-footer
    :mask-closable="false"
    @close="handleClose"
    @confirm="handleSubmit"
  >
    <tiny-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-position="top"
      size="large"
    >
      <tiny-row :gutter="20">
        <tiny-col :span="24">
          <tiny-form-item label="图片上传">
            <tiny-file-upload
              v-model="formData.image_url"
              :action="imageUploadUrl"
              :max-size="1024 * 1024 * 5"
              :accept="'.jpg,.jpeg,.png'"
              class="upload-demo picture-demo"
              :auto-upload="true"
              :limit="1"
              is-hidden
              list-type="picture-card"
              :headers="headers"
              @success="handleUploadSuccess"
            >
              <template #default>
                <tiny-icon-plus class="tiny-svg-size" />
              </template>
              <template #file="{ file }">
                <div>
                  <img
                    class="tiny-upload-list__item-thumbnail"
                    :src="file.url"
                    alt=""
                  />
                  <span class="tiny-upload-list__item-actions">
                    <span
                      v-if="!disabled"
                      class="tiny-upload-list__item-delete"
                      @click="removePicture(file)"
                    >
                      <tiny-icon-del class="tiny-svg-size" />
                    </span>
                  </span>
                </div>
              </template>
            </tiny-file-upload>
          </tiny-form-item>
        </tiny-col>
        <tiny-col :span="24">
          <tiny-form-item label="轮播图标题" prop="title">
            <tiny-input
              v-model="formData.title"
              placeholder="请输入轮播图标题"
              :maxlength="255"
              show-word-limit
            ></tiny-input>
          </tiny-form-item>
        </tiny-col>
        <tiny-col :span="24">
          <tiny-form-item label="图片预览">
            <div class="image-preview-container">
              <img
                v-if="formData.image_url"
                :src="formData.image_url"
                alt="轮播图预览"
                class="preview-image"
              />
              <div v-else class="no-image"> 暂无图片预览 </div>
            </div>
          </tiny-form-item>
        </tiny-col>
        <tiny-col :span="24">
          <tiny-form-item label="排序值" prop="sort">
            <tiny-input
              v-model="formData.sort"
              :min="0"
              :max="9999"
              placeholder="请输入排序值"
            ></tiny-input>
          </tiny-form-item>
        </tiny-col>
        <tiny-col :span="24">
          <tiny-form-item label="是否可见" prop="is_visible">
            <tiny-switch v-model="formData.is_visible"></tiny-switch>
          </tiny-form-item>
        </tiny-col>
      </tiny-row>
    </tiny-form>
    <template #footer>
      <tiny-button @click="handleClose">取消</tiny-button>
      <tiny-button type="info" @click="handleSubmit">
        {{ isEdit ? '更新' : '创建' }}
      </tiny-button>
    </template>
  </tiny-drawer>
</template>

<script lang="ts" setup>
  import {
    ref,
    reactive,
    watch,
    toRaw,
    toRef,
    onMounted,
    defineEmits,
    defineProps,
  } from 'vue';
  import {
    Drawer as TinyDrawer,
    Form as TinyForm,
    FormItem as TinyFormItem,
    Input as TinyInput,
    Button as TinyButton,
    Row as TinyRow,
    Col as TinyCol,
    Switch as TinySwitch,
    Modal,
    TinyFileUpload,
  } from '@opentiny/vue';
  import type { BannerForm, BannerInfo } from '@/api/banner';
  import { addBanner, updateBanner } from '@/api/banner';
  import {
    iconPlus,
    iconView,
    iconDownload,
    iconDel,
  } from '@opentiny/vue-icon';
  import { getToken } from '@/utils/auth';

  const baseUrl = import.meta.env.VITE_BASE_API || '';
  const imageUploadUrl = `${baseUrl}/api/upload`;
  const TinyIconPlus = iconPlus();
  const TinyIconView = iconView();
  const TinyIconDownload = iconDownload();
  const TinyIconDel = iconDel();
  const pictureUploadRef = ref();

  const visible = ref<boolean>(false);
  const props = defineProps<{
    visible: boolean;
    bannerData?: BannerInfo;
  }>();

  const emit = defineEmits(['update:visible', 'success']);
  const propsVisible = toRef(props, 'visible');
  const propsBannerData = toRef(props, 'bannerData');

  const formRef = ref();
  const isEdit = ref(false);
  const formData = reactive<BannerForm>({
    id: undefined,
    title: '',
    image_url: '',
    is_visible: true,
    sort: 0,
  });

  const rules = {
    title: [
      {
        required: true,
        message: '请输入轮播图标题',
        trigger: 'blur',
      },
    ],
    image_url: [
      {
        required: true,
        message: '请输入图片URL',
        trigger: 'blur',
      },
      {
        type: 'url',
        message: '请输入有效的图片URL',
        trigger: 'blur',
      },
    ],
    sort: [
      {
        required: true,
        message: '请输入排序值',
        trigger: 'blur',
      },
      {
        type: 'number',
        message: '排序值必须为数字',
        trigger: 'blur',
      },
    ],
  };

  const headers = ref({
    Authorization: `${getToken()}`,
    token: `${getToken()}`,
  });

  function previewPicture(file) {
    // dialogImageUrl.value = file.url
    // dialogVisible.value = true
  }

  function removePicture(file) {
    pictureUploadRef.value.handleRemove(file);
  }

  const handleUploadSuccess = (res: any, file: any) => {
    console.log('res:', res);
    if (res.code === 200) {
      file.url = res.url;
      formData.image_url = res.url;
      formData.title = res.orgName;
    } else {
      Modal.message({
        message: res.msg || '图片上传失败',
        status: 'error',
      });
    }
  };

  // 重置表单
  const resetForm = () => {
    formData.id = undefined;
    formData.title = '';
    formData.image_url = '';
    formData.is_visible = true;
    formData.sort = 0;
    formRef.value?.resetFields();
  };

  // 关闭抽屉
  const handleClose = () => {
    visible.value = false;
    emit('update:visible', false);
    setTimeout(() => {
      resetForm();
    }, 300);
  };

  // 提交表单
  const handleSubmit = () => {
    formRef.value?.validate(async (valid: boolean) => {
      if (valid) {
        try {
          const submitData = { ...toRaw(formData) };
          let res;

          if (isEdit.value) {
            res = await updateBanner(submitData);
          } else {
            res = await addBanner(submitData);
          }

          if (res.code === 200) {
            Modal.message({
              message: res.msg,
              status: 'success',
            });
            emit('success');
            handleClose();
          } else {
            Modal.message({
              message: res.msg,
              status: 'error',
            });
          }
        } catch (error: any) {
          Modal.message({
            message: error.message || '操作失败',
            status: 'error',
          });
        }
      }
    });
  };

  // 监听数据变化
  watch(
    () => propsBannerData.value,
    (val) => {
      if (val && val.id) {
        isEdit.value = true;
        // 转换后端数据到表单格式
        formData.id = val.id;
        formData.title = val.title;
        formData.image_url = val.image_url;
        formData.is_visible = val.is_visible;
        formData.sort = val.sort;
      } else {
        isEdit.value = false;
        resetForm();
      }
    },
    { immediate: true, deep: true },
  );
  watch(visible, (newValue, oldValue) => {
    emit('update:visible', visible.value);
  });
  watch(propsVisible, (newValue, oldValue) => {
    visible.value = propsVisible.value;
    console.log('watch propsVisible.value', propsVisible.value);
  });

  onMounted(() => {
    visible.value = propsVisible.value;
  });
</script>

<style scoped>
  .btn-group {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }

  .btn-group :deep(.tiny-btn) {
    margin-left: 10px;
  }

  .image-preview-container {
    width: 100%;
    height: 200px;
    border: 1px solid #e8e8e8;
    border-radius: 4px;
    display: flex;
    justify-content: center;
    align-items: center;
    overflow: hidden;
    background-color: #f5f5f5;
  }

  .preview-image {
    max-width: 100%;
    max-height: 100%;
    object-fit: contain;
  }

  .no-image {
    color: #999;
    font-size: 14px;
  }
</style>

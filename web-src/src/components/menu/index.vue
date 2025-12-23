<template>
  <div class="menu-router">
    <tiny-tree-menu
      ref="tree"
      :data="MenuData"
      :show-filter="false"
      node-key="id"
      wrap
      :default-expanded-keys="expandeArr"
      only-check-children
      check-strictly
      @current-change="currentChange"
    >
      <template #default="slotScope">
        <template v-for="(item, index) in routerTitle" :key="index">
          <span v-if="slotScope.label === item.label" class="menu-title">
            <component :is="item.customIcon"></component>
            <span>{{ $t(item.locale) }}</span>
          </span>
        </template>
      </template>
    </tiny-tree-menu>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, watch, computed, unref } from "vue";
import { RouteRecordNormalized, useRouter } from "vue-router";
import { TreeMenu as tinyTreeMenu } from "@opentiny/vue";
import router from "@/router";
import * as icons from "@opentiny/vue-icon";
import { useTabStore } from "@/store";
import { ITreeNodeData } from "@/router/guard/index";
// import { useDeepClone } from "@/hooks/useDeepClone";


// import { useMenuStore } from "@/store/modules/router";
// const menuStore = useMenuStore();
// await menuStore.getMenuList();
// const rawMenuData = computed(() => useDeepClone(unref(menuStore.menuList)));

type SideMenuData = (ITreeNodeData & { meta: { url: string } })[];

let routerTitle = [] as any;

const currRole = ref('');

const tree = ref();
const expandeArr = ref<(string | number)[]>([]);
const tabStore = useTabStore();
const routerUse = useRouter();
// 获取路由数据
const appRoute = computed(() => {
  return router
    .getRoutes()
    .find((el: { name: string }) => el.name === "root") as RouteRecordNormalized;
});
const copyRouter = JSON.parse(JSON.stringify(appRoute.value.children));
copyRouter.sort((a: RouteRecordNormalized, b: RouteRecordNormalized) => {
  return (a.meta.order || 0) - (b.meta.order || 0);
});
let treeDataRef = ref(copyRouter);
const treeDataForEach = (arr: any[]) => {
  return arr.filter((e: { children: any[]; meta: { hideInMenu: any } }) => {
    if (e.children) {
      e.children = e.children.filter((v: { meta: { hideInMenu: any } }) => {
        return !v.meta.hideInMenu;
      });
      treeDataForEach(e.children);
    }
    return !e.meta.hideInMenu;
  });
};

const filtter = (treeNodeDatas: ITreeNodeData[]) => {
  const menus: SideMenuData = [];
  for (let i = 0; i < treeNodeDatas.length; i += 1) {
    const treeNodeData = treeNodeDatas[i];
    // 判断角色
    // console.log("判断角色:",treeNodeData.meta.roles.includes(currRole),currRole.value);
    if (!treeNodeData.meta.roles.includes(currRole.value)){
      continue;
    }
    let url = "";
    if (treeNodeData.url) {
      url = treeNodeData.url ? treeNodeData.url! : "";
      delete treeNodeData.url;
    }
    const temp = {} as any;
    temp.label = treeNodeData.label;
    temp.locale =
      treeNodeData.meta && treeNodeData.meta.locale
        ? treeNodeData.meta.locale
        : treeNodeData.locale;
    if (treeNodeData.customIcon) {
      temp.customIcon = icons[treeNodeData.customIcon]();
    }
    routerTitle.push(temp);
    menus.push({
      ...treeNodeData,
      meta: {
        url,
      },
      children: [...filtter(treeNodeData.children ?? [])],
    });
  }
  return menus;
};

const MenuData = computed(() => {
  if (routerTitle.length) {
    routerTitle = [];
  }
  return filtter(treeDataForEach(treeDataRef.value));
});

const currentChange = (data: any, node) => {
  if (!node.isLeaf) {
    return;
  }
  router.replace({ name: data.label });
};

const findId = (name: string, path: string) => {
  const dfs = (item, url: string[]) => {
    if (url.join("/") === path) {
      return item.id;
    }
    const len = item.children.length ?? 0;
    for (let i = 0; i < len; i += 1) {
      if (item.children?.[i]) {
        const id = dfs(
          item.children[i],
          [...url, item.children[i].meta.url].filter((p) => p.length)
        );
        if (id !== undefined) {
          return id;
        }
      }
    }
    return undefined;
  };
  for (let i = 0; i < MenuData.value.length; i += 1) {
    const menu = MenuData.value[i];
    const data = dfs(menu, [
      import.meta.env.VITE_CONTEXT.replace(/\/$/, ""),
      menu.meta.url.replace(/\/$/, ""),
    ]);
    if (data !== undefined) {
      return data;
    }
  }
  return -1;
};

onMounted(() => {
  watch(
    () => tabStore.current,
    () => {
      if (!tabStore.current) {
        return;
      }
      const key = findId(tabStore.current.name, tabStore.current.link);
      tree.value.setCurrentKey(key);

      let parentId:any = null;
      parentId =
        routerUse.currentRoute.value.matched.length >= 2
          ? routerUse.currentRoute.value.matched[1].name
          : routerUse.currentRoute.value.matched[0].name;
      if (parentId && !expandeArr.value.includes(parentId)) {
        expandeArr.value = expandeArr.value.concat(parentId);
      }
    },
    { deep: true, immediate: true }
  );
  currRole.value = localStorage.getItem("USER_ROLE_STORE_STATE") || ""
});
</script>

<style scoped></style>

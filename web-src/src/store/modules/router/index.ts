import { defineStore } from 'pinia';
import { getRoleMenu } from '@/api/menu';

export const useMenuStore = defineStore('menu', {
  state() {
    return {
      menuList: [] as any[],
      flatMenuList: [] as any[],
    };
  },
  actions: {
    async getMenuList() {
      
      const { data } = await getRoleMenu();
      this.menuList = data;
      this.menuListFlat();
      return data;
    },
    menuListFlat() {
      this.flatMenuList = [];
      const dfs = (item: any) => {
        this.flatMenuList.push(item);
        for (let i = 0; i < item.children.length; i += 1) {
          dfs(item.children[i]);
        }
      };
      for (let i = 0; i < this.menuList.length; i += 1) {
        dfs(this.menuList[i]);
      }
    },
  },
});

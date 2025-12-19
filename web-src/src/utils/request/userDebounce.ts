import { debounce } from '@/utils/request/debounce'
import app from '@/App.vue';

app.component('Debounce', {
  abstract: true,
  props: ['time', 'events', 'immediate'],
  created() {
    this.eventKeys = this.events && this.events.split(',');
  },
  render() {
    const vnode = this.$slots.default()[0];
    this.eventKeys.forEach((key:any) => {
      vnode.props[`on${key}`] = debounce(
        vnode.props[`on${key}`],
        this.time,
        vnode,
        this.immediate,
      );
    });
    return vnode;
  },
});

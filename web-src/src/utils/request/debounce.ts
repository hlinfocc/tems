export const debounce = (func:any, time:any, ctx:any, immediate:any) => {
  let timer:any;
  // @ts-ignore
  const rtn = (...params) => {
    clearTimeout(timer);
    if (immediate) {
      let callNow = !timer;
      timer = setTimeout(() => {
        timer = null;
      }, time);
      if (callNow) func.apply(ctx, params);
    } else {
      timer = setTimeout(() => {
        func.apply(ctx, params);
      }, time);
    }
  };
  return rtn;
};

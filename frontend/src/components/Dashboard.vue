<script setup>
import {reactive} from 'vue'
import {Greet} from '../../wailsjs/go/main/App'
import * as rt from "../../wailsjs/runtime/runtime.js"; // the runtime for Wails2

const data = reactive({
  state: 0,
  path: "",
  total: 0,
  success: 0,
  fail: 0
})


rt.EventsOn("ProcessStart", async (msg) => {
    data.state = 1;
    data.path = msg.path;
    data.total = msg.total;
});

rt.EventsOn("ProcessUpdate", async (msg) => {
    data.success = msg.success;
    data.fail = msg.fail;
});

rt.EventsOn("ProcessDone", async (msg) => {
    data.state = 2;
});

</script>

<template>
  <main>
    <div v-if="data.state==0" class="tip">👆🏻 左上角菜单选择 &lt;文件转换&gt; 或 &lt;文件夹转换&gt; 开始！</div>
    <div v-if="data.state==1" class="tip">转换中...请稍后...</div>
    <div v-if="data.state==2" class="tip">转换完成！</div>
    <div v-if="data.state > 0">
      <div class="result">路径：{{ data.path }}</div>
      <div class="result">总共：{{ data.total }}</div>
      <div class="result">成功：{{ data.success }}</div>
      <div class="result">失败：{{ data.fail }}</div>
    </div>
  </main>
</template>

<style scoped>


.tip {
  height: 20px;
  line-height: 20px;
  margin: 1.5rem auto;
}
.result {
  /* height: 20px;
  line-height: 20px;
  margin: 1.5rem auto; */
  color: darkgreen;
}
</style>

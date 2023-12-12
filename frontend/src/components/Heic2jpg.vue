
<template>
  <main>

    <!-- 上传按钮 -->
    <div class="clearfix">
      <a-row>
        <a-col :span="12">
          <a-upload
            :multiple="true"
            accept=".HEIC,.heic"
            :file-list="fileList"
            :before-upload="beforeUpload"
            @remove="handleRemove"
          >
            <a-button>
              <upload-outlined></upload-outlined>
              选择文件
            </a-button>

            <!-- <template #itemRender="{ file, actions }">
          <a-space>
            <span :style="file.status === 'error' ? 'color: red' : ''">{{
              file.name
            }}</span>
            <a href="javascript:;" @click="actions.convert">转换</a>
            <a href="javascript:;" @click="actions.remove">删除</a>
          </a-space>
        </template> -->
          </a-upload>
        </a-col>
        <a-col :span="12">
            
    <div class="clearfix">
         <a-tooltip placement="top">
        <template #title>
          <span>cmd/ctrl+s 进行设置</span>
        </template>
        保存路径：<a-tag> {{ title }} </a-tag>
      </a-tooltip>
    </div>
    
          <a-button
            type="primary"
            :disabled="fileList.length === 0"
            v-if="fileList.length > 0"
            :loading="uploading"
            style="margin-top: 16px;"
            @click="handleUpload"
          >
            {{ uploading ? "转换中" : "开始转换" }}
          </a-button>
        </a-col>
      </a-row>
    </div>
  </main>
</template>
<script setup>
import { onMounted, ref, nextTick } from "vue";
import request from "umi-request";
import { message } from "ant-design-vue";
import * as App from "../../wailsjs/go/main/App.js";
import * as rt from "../../wailsjs/runtime/runtime.js"; // the runtime for Wails2

import { notification } from 'ant-design-vue';


const initLoading = ref(false);
const loading = ref(false);
const data = ref([]);
const title = ref('');
const savepath = ref('');

const fileList = ref([]);
const uploading = ref(false);
const handleRemove = (file) => {
  const index = fileList.value.indexOf(file);
  const newFileList = fileList.value.slice();
  newFileList.splice(index, 1);
  fileList.value = newFileList;
};
const beforeUpload = (file) => {
  const files = [...(fileList.value || []), file];
  fileList.value = files;
  return false;
};

const handleUpload = async() => {
    // uploading.value = true;
    const files = [...(fileList.value || [])];
    notification['info']({
        message: 'Notification Title',
        description:
        'handleUpload.' + files.length,
    });
    for (let index = 0; index < files.length; index++) {
        const element = files[index];
          var filepath = URL.createObjectURL(element);

         notification['info']({
                message: 'Notification Title',
                description:
                'element.' + filepath,
            });
        
        let result  = await App.Heic2jpg(element)
        if(result == ''){
            console.log('成功', element)
        }else{
            console.log('失败', element)
        }
    }
    uploading.value = false;
};



const savePathSelect = () => {
    window.runtime
}

rt.EventsOn("SettingSavePath", async (msg) => {
    savepath.value = msg
    title.value = msg;
});
</script>
<style scoped>
.demo-loadmore-list {
  min-height: 350px;
}
</style>
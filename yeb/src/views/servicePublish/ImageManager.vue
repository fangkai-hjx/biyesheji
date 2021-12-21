<template>
  <div class="mod-config">
    <el-form
      :inline="true"
      :model="dataForm"
      @keyup.enter.native="getDataList()"
    >
      <el-form-item>
        <el-select
          v-model="value"
          placeholder="请选择工作空间"
          @change="getImageProjects"
        >
          <el-option
            v-for="item in options"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          >
          </el-option>
        </el-select>
        工作空间
      </el-form-item>
      <el-form-item>
        <el-input
          v-model="dataForm.key"
          placeholder="亲输入镜像名称"
          clearable
        ></el-input>
      </el-form-item>
      <el-form-item>
        <el-button @click="getDataList()">查询</el-button>
        <!-- <el-button v-if="isAuth('member:member:save')" type="primary" @click="addOrUpdateHandle()">新增</el-button>
        <el-button v-if="isAuth('member:member:delete')" type="danger" @click="deleteHandle()" :disabled="dataListSelections.length <= 0">批量删除</el-button>-->
      </el-form-item>
    </el-form>
    <el-table
      :data="dataList"
      border
      v-loading="dataListLoading"
      @selection-change="selectionChangeHandle"
      style="width: 100%"
    >
      <el-table-column
        type="selection"
        header-align="center"
        align="center"
        width="50"
      ></el-table-column>
      <el-table-column
        prop="project_id"
        header-align="center"
        align="center"
        label="ID"
      ></el-table-column>
      <el-table-column
        prop="image_name"
        header-align="center"
        align="center"
        label="镜像名"
      ></el-table-column>
      <el-table-column
        prop="tag"
        header-align="center"
        align="center"
        label="标签"
      ></el-table-column>
      <el-table-column
        prop="pull_count"
        header-align="center"
        align="center"
        label="使用次数"
      ></el-table-column>
       <el-table-column
        prop="create_time"
        header-align="center"
        align="center"
        label="创建时间"
      ></el-table-column>
       <el-table-column
        prop="update_time"
        header-align="center"
        align="center"
        label="更新时间"
      ></el-table-column>
      <el-table-column
        prop="status"
        header-align="center"
        align="center"
        label="启用状态"
      >
        <template slot-scope="scope">
          <el-switch
            v-model="scope.row.status"
            active-color="#13ce66"
            inactive-color="#ff4949"
            :active-value="1"
            :inactive-value="0"
          ></el-switch>
        </template>
      </el-table-column>     
       <el-table-column fixed="right" header-align="center" align="center" label="操作">
              <template slot-scope="scope">
                <el-button type="text" size="small" @click="relationRemove(scope.row.attrId)">删除</el-button>
              </template>
            </el-table-column>
    </el-table>
    <el-pagination
      @size-change="sizeChangeHandle"
      @current-change="currentChangeHandle"
      :current-page="pageIndex"
      :page-sizes="[10, 20, 50, 100]"
      :page-size="pageSize"
      :total="totalPage"
      layout="total, sizes, prev, pager, next, jumper"
    ></el-pagination>
    <!-- 弹窗, 新增 / 修改 -->
    <add-or-update
      v-if="addOrUpdateVisible"
      ref="addOrUpdate"
      @refreshDataList="getDataList"
    ></add-or-update>
  </div>
</template>

<script>
import Axios from "axios";
export default {
  data() {
    return {
      value: "",
      options: [],
      dataForm: {
        key: "",
      },
      dataList: [],
      pageIndex: 1,
      pageSize: 10,
      totalPage: 0,
      dataListLoading: false,
      dataListSelections: [],
      addOrUpdateVisible: false,
    };
  },
  mounted() {
    this.getDataList();
    this.getImageProjects();
  },
  methods: {
    clear() {
      this.options = [];
      this.dataList = [];
    },
    getImageForProject() {
      this.tableData = [];
      var api =
        "http://202.38.247.217:8080/api/v1/rest/img/workspace/" +
        this.value +
        "/img";
      Axios.get(api)
        .then((res) => {
          for (const key in res.data.data) {
            this.dataList.push({
              project_id: res.data.data[key]["project_id"],
              image_name: res.data.data[key]["image_name"],
              tag: res.data.data[key]["tag"],
              pull_count: res.data.data[key]["pull_count"],
              create_time: res.data.data[key]["create_time"],
              update_time: res.data.data[key]["update_time"],
            });            
          }
          this.dataListLoading = false;
        })
        .catch((error) => {});
    },
    getImageProjects() {
      this.clear();
      var api = "http://202.38.247.217:8080/api/v1/rest/img/workspace/project";
      Axios.get(api)
        .then((res) => {
          console.log(res.data.data);
          for (const key in res.data.data) {
            this.options.push({
              value: res.data.data[key],
              label: res.data.data[key],
            });
          }
          this.dataListLoading = false;
        })
        .catch((error) => {});
      this.getImageForProject();
    },
    // 获取数据列表,
    getDataList() {      
      this.dataListLoading = true;
      var api = "http://202.38.247.217:8080/api/v1/rest/img/workspace/all";
      Axios.get(api)
        .then((res) => {
          var s = res.data.data;
          for (const key in s) {
            console.log(s["Images"]);
            for (const key in s["Images"]) {
              console.log(s["Images"]["image_name"]);
            }
          }
          this.dataListLoading = false;
        })
        .catch((error) => {});
    },
    // 每页数
    sizeChangeHandle(val) {
      this.pageSize = val;
      this.pageIndex = 1;
      this.getDataList();
    },
    // 当前页
    currentChangeHandle(val) {
      this.pageIndex = val;
      this.getDataList();
    },
    // 多选
    selectionChangeHandle(val) {
      this.dataListSelections = val;
    },
    // 新增 / 修改
    addOrUpdateHandle(id) {
      this.addOrUpdateVisible = true;
      this.$nextTick(() => {
        this.$refs.addOrUpdate.init(id);
      });
    },
    // 删除
    //   deleteHandle (id) {
    //     var ids = id ? [id] : this.dataListSelections.map(item => {
    //       return item.id
    //     })
    //     this.$confirm(`确定对[id=${ids.join(',')}]进行[${id ? '删除' : '批量删除'}]操作?`, '提示', {
    //       confirmButtonText: '确定',
    //       cancelButtonText: '取消',
    //       type: 'warning'
    //     }).then(() => {
    //       this.$http({
    //         url: this.$http.adornUrl('/member/member/delete'),
    //         method: 'post',
    //         data: this.$http.adornData(ids, false)
    //       }).then(({data}) => {
    //         if (data && data.code === 0) {
    //           this.$message({
    //             message: '操作成功',
    //             type: 'success',
    //             duration: 1500,
    //             onClose: () => {
    //               this.getDataList()
    //             }
    //           })
    //         } else {
    //           this.$message.error(data.msg)
    //         }
    //       })
    //     })
    //   }
  },
};
</script>
